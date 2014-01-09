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
        $.ajax($("#url").val() + $('#query-input').val(), {
            contentType: "application/json",
            type: "GET",
            success: function (data) {
                $('#query-result').html(JSON.stringify(data.Data));
            },
            error: function (data) {
                $('#query-result').html(data.status + " " + data.statusText);
            }
        })
    });

    $('#start-upload').click(function (e) {
        sendElementsToServer(uploadData.objectsToUpload);
    });

    $('#post-button').click(function (e) {
        $.ajax($("#url").val() + '/data', {
            contentType: "application/json",
            type: "GET",
            data: $('#post-data').val(),
            type: "POST",
            complete: function (data) {
                $('#post-result').html(data.status + " " + data.statusText);
            }
        })
    });

    $('#histogram-button').click(function (e) {
        var chart = drawHistogram([3, 2, 1, 6, 10, 5, 13, 9, 14, 21, 23, 66, 47, 14, 5, 2],
            ['> 48.00 =< 51.81', '> 51.81 =< 54.63', '> 54.63 =< 57.44', '> 57.44 =< 60.25', '> 60.25 =< 63.06', '> 63.06 =< 65.88', '> 65.88 =< 68.69', '> 68.69 =< 71.50', '> 71.50 =< 74.31',
            '> 74.31 =< 77.13', '> 77.13 =< 79.94', '> 79.94 =< 82.75', '> 82.75 =< 85.56', '> 85.56 =< 88.38', '> 88.38 =< 91.19', '> 91.19 =< 94.00'],
            'chart', 'zajebisty histogram');
    });
});

var elementToJSONString = function (element) {
    return JSON.stringify({
        "tag": parseInt(element.tag, 10),
        "metric": parseInt(element.metric, 10),
        "time": parseInt(element.time, 10),
        "value": parseFloat(element.value)
    });
};

var sendElementsToServer = function (dataToSend) {
    dataToSend.forEach(function (dataElement) {
        $.ajax($("#url").val() + "/data", {
            contentType: "application/json",
            dataType: "json",
            data: elementToJSONString(dataElement),
            type: "POST",
            success: function () {

            },
            error: function () {
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
                    fontSize: '8px',
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


