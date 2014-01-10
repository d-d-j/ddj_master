var uploadData;

$(document).ready(function () {
    $('#post-data').val("{\n\"tag\":1,\n\"metric\":2,\n\"time\":1383501407,\n\"value\":0.5\n}");

    $("#file-input").change(function (e) {
        var file = e.target.files[0];
        var reader = new FileReader();
        reader.readAsText(file);
        reader.onload = function (event) {
            var csv = event.target.result;
            uploadData = {objectsToUpload: $.csv.toObjects(csv)};

            var dataDiv = $('#data-to-upload');
            dataDiv.html('');
            uploadData.objectsToUpload.forEach(function (element) {
                dataDiv.append(JSON.stringify(element));
            });
        };
    });

    $('#query-button').click(function (e) {
        $('#query-result').html(" ");
        $.ajax($("#url").val() + $('#query-input').val(), {
            contentType: "application/json",
            type: "GET",
            success: function (data) {
                $('#query-result').html(syntaxHighlight(JSON.stringify(data.Data, undefined, 4)));
            },
            error: function (data) {
                $('#query-result').html('Response: <span class="error-result">' + data.status + " " + data.statusText + '</span>');
            }
        })
    });

    $('#start-upload').click(function (e) {
        if (uploadData != undefined && uploadData.objectsToUpload.length > 0) {
            sendElementsToServer(uploadData.objectsToUpload);
        } else {
            $('#upload-result').html('<span class="error-result"> No data to upload </span>');
        }
    });

    $('#post-button').click(function (e) {
        $('#post-result').html(" ");
        $.ajax($("#url").val() + '/data', {
            contentType: "application/json",
            data: $('#post-data').val(),
            type: "POST",
            success: function (data) {
                $('#post-result').html('Response: <span class="success-result">' + data + '</span>');
            },
            error: function (data) {
                $('#post-result').html('Response: <span class="error-result">' + data.status + " " + data.statusText + '</span>');
            }
        })
    });

    $('#histogram-button').click(function (e) {
        var tags = $('#histogram-tags').val();
        var metrics = $('#histogram-metrics').val();
        var timeFrom = $('#histogram-time-from').val();
        var timeTo = $('#histogram-time-to').val();
        var numBuckets = $('#histogram-buckets').val();
        var histogramOption = $('#histogram-option').val();
//        http://localhost:8888/data/metric/all/tag/all/time/194-317/aggregation/histogramByTime/from/194/to/317/buckets/1
        var queryUrl = "/data/metric/" + metrics + "/tag/" + tags + "/time/all/aggregation/histogramBy" + histogramOption + "/from/" + timeFrom + "/to/" + timeTo + "/buckets/" + numBuckets;

        var bucketSize = (parseFloat(timeTo) - parseFloat(timeFrom)) / parseFloat(numBuckets);

        var buckets = [];

        for (var i = 0; i < numBuckets; i++) {
            var begin = Math.round((parseFloat(timeFrom) + parseFloat(i * bucketSize)) * 100) / 100;
            var end = Math.round((parseFloat(timeFrom) + parseFloat(((i + 1) * bucketSize))) * 100) / 100;
            buckets.push('[ ' + String(begin) + ', ' + String(end) + ')');
        }
        $.ajax($("#url").val() + queryUrl, {
            contentType: "application/json",
            type: "GET",
            success: function (data) {
                console.log(data.Data.length);
                if (data.Data.length > 0) {
                    var chart = drawHistogram(data.Data[0].Data, buckets, 'histogram-chart', "histogram by " + histogramOption);
                } else {
                    $('#histogram-chart').html("No data to display");
                }
            },
            error: function (data) {
                $('#histogram-chart').html('Response: <span class="error-result">' + data.status + " " + data.statusText + '</span>');
            }
        })
    });

    $('#series-button').click(function (e) {
        var tag1 = $('#series-tag1').val();
        var tag2 = $('#series-tag2').val();
        var metric = $('#series-metric').val();
        var timeFrom = $('#series-time-from').val();
        var timeTo = $('#series-time-to').val();
        var seriesOption = $('#series-option').val();
        var numSamples = $('#series-num-samples').val();
        <!--/metric/%s/tag/%s/time/from/%d/to/%d/aggregation/series/sum/samples/1-->
        var queryUrl = "/data/metric/" + metric + "/tag/" + tag1 + ',' + tag2 + "/time" + "/from/" + timeFrom + "/to/" + timeTo + "/aggregation/series/" + seriesOption + "/samples/" + numSamples;
        var series1Url = "/data/metric/" + metric + "/tag/" + tag1 + "/time/" + timeFrom + "-" + timeTo + "/aggregation/none";
        var series2Url = "/data/metric/" + metric + "/tag/" + tag2 + "/time/" + timeFrom + "-" + timeTo + "/aggregation/none";

        $.when(
                $.ajax($("#url").val() + queryUrl, {
                    contentType: "application/json",
                    type: "GET"
                }), $.ajax($("#url").val() + series1Url, {
                    contentType: "application/json",
                    type: "GET"
                }), $.ajax($("#url").val() + series2Url, {
                    contentType: "application/json",
                    type: "GET"
                })
            ).then(function (aggregatedData, series1Data, series2Data) {
                drawSeriesAggregation(parseElementsToPoints(series1Data[0].Data), parseElementsToPoints(series2Data[0].Data),
                    parseSeriesAggregationResultsToPoints(aggregatedData[0].Data[0].Data, timeFrom, timeTo, numSamples), 'series-chart', 'aggregation');

            }, function () {
            });
    });
});

var parseSeriesAggregationResultsToPoints = function (results, timeFrom, timeTo, samples) {
    var gap = (timeTo - timeFrom) / samples;


    var points = [];
    for (var i = 0; i < results.length; i++) {
        points.push([timeFrom + i * gap, results[i] ]);
    }
    return points;
};

var parseElementsToPoints = function (elements) {
    var points = [];

    elements.forEach(function (element) {
        points.push([element.Time, element.Value]);
    });

    return points;
};

var elementToJSONString = function (element) {
    return JSON.stringify({
        "tag": parseInt(element.tag, 10),
        "metric": parseInt(element.metric, 10),
        "time": parseInt(element.time, 10),
        "value": parseFloat(element.value)
    });
};

var sendElementsToServer = function (dataToSend) {
    $('#upload-result').html(" ");
    dataToSend.forEach(function (dataElement) {
        $.ajax($("#url").val() + "/data", {
            contentType: "application/json",
            dataType: "json",
            data: elementToJSONString(dataElement),
            type: "POST",
            error: function (data) {
                if (data.status === 202) {
                    $('#upload-result').html('Response: <span class="success-result">' + data.status + " " + data.statusText + '</span>');
                } else {
                    $('#upload-result').html('Response: <span class="error-result">' + data.status + " " + data.statusText + '</span>');

                }
            }
        })
    })
};


var drawHistogram = function (series, bins, divId, title) {
    return new Highcharts.Chart({
        chart: {
            renderTo: divId,
            defaultSeriesType: 'column',
            borderWidth: 0,
            backgroundColor: '#eee',
            borderWidth: 1,
            borderColor: '#ccc',
            plotBackgroundColor: '#fff',
            plotBorderWidth: 1,
            plotBorderColor: '#ccc'
        },
        credits: {enabled: false},
        exporting: {enabled: false},
        title: {text: title},
        legend: {
            //enabled:false
        },
        tooltip: {
            borderWidth: 1,
            formatter: function () {
                return '<b>Range:</b><br/> ' + this.x + '<br/>' +
                    '<b>Count:</b> ' + this.y;
            }
        },
        plotOptions: {
            column: {
                shadow: false,
                borderWidth: .5,
                borderColor: '#666',
                pointPadding: 0,
                groupPadding: 0,
                color: 'rgba(204,204,204,.85)'
            },
            spline: {
                shadow: false,
                marker: {
                    radius: 1
                }
            },
            areaspline: {
                color: 'rgb(69, 114, 167)',
                fillColor: 'rgba(69, 114, 167,.25)',
                shadow: false,
                marker: {
                    radius: 1
                }
            }
        },
        xAxis: {
            categories: bins,
            labels: {
                rotation: -90,
                y: 40,
                style: {
                    fontSize: '12px',
                    fontWeight: 'normal',
                    color: '#333'
                }
            },
            lineWidth: 0,
            lineColor: '#999',
            tickLength: 70,
            tickColor: '#ccc'
        },
        yAxis: {
            title: {text: ''},
            //maxPadding:0,
            gridLineColor: '#e9e9e9',
            tickWidth: 1,
            tickLength: 3,
            tickColor: '#ccc',
            lineColor: '#ccc',
            tickInterval: 25
            //endOnTick:false,
        },
        series: [
            {
                name: 'Bins',
                data: series
            }
        ]
    });

};

var drawSeriesAggregation = function (series1, series2, aggregated, divId, title) {
    console.log(series1);
    return new Highcharts.Chart({
        chart: {
            renderTo: divId,
            borderWidth: 0,
            backgroundColor: '#eee',
            borderWidth: 1,
            borderColor: '#ccc',
            plotBackgroundColor: '#fff',
            plotBorderWidth: 1,
            plotBorderColor: '#ccc'
        },
        credits: {enabled: false},
        exporting: {enabled: false},
        title: {text: title},
        legend: {
            //enabled:false
        },
        tooltip: {
            borderWidth: 1,
            formatter: function () {
                return '<b>X:</b><br/> ' + this.x + '<br/>' +
                    '<b>Y:</b> ' + this.y;
            }
        },
        plotOptions: {

            spline: {
                shadow: false,
                marker: {
                    radius: 1
                }
            }
        },
        xAxis: {
            labels: {
                rotation: -90,
                y: 40,
                style: {
                    fontSize: '12px',
                    fontWeight: 'normal',
                    color: '#333'
                }
            },
            lineWidth: 0,
            lineColor: '#999',
            tickLength: 10,
            tickColor: '#ccc'
        },
        yAxis: {
            title: {text: ''},
            //maxPadding:0,
            gridLineColor: '#e9e9e9',
            tickWidth: 1,
            tickLength: 3,
            tickColor: '#ccc',
            lineColor: '#ccc',
            tickInterval: 0.1
            //endOnTick:false,
        },
        series: [
            {
                name: 'Series1',
                type: 'spline',
                data: series1
            } ,
            {
                name: 'Series2',
                type: 'spline',
                data: series2
            },
            {
                name: 'sum',
                type: 'spline',
                data: aggregated
            }
        ]
    });

};


var syntaxHighlight = function (json) {
    json = json.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;');
    return json.replace(/("(\\u[a-zA-Z0-9]{4}|\\[^u]|[^\\"])*"(\s*:)?|\b(true|false|null)\b|-?\d+(?:\.\d*)?(?:[eE][+\-]?\d+)?)/g, function (match) {
        var cls = 'number';
        if (/^"/.test(match)) {
            if (/:$/.test(match)) {
                cls = 'key';
            } else {
                cls = 'string';
            }
        } else if (/true|false/.test(match)) {
            cls = 'boolean';
        } else if (/null/.test(match)) {
            cls = 'null';
        }
        return '<span class="' + cls + '">' + match + '</span>';
    });
}