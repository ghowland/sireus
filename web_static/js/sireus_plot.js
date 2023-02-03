function CreatePlot(element_id, title, plot_x, plot_y, plot_selected_x, plot_selected_y)
{
    var trace1 = {
        x: plot_x,
        y: plot_y,
        // mode: 'lines+markers',
        mode: 'lines',
        name: 'Curve',
        line: {shape: 'linear'},
        type: 'scatter'
    };

    var data = [trace1]

    if (plot_selected_x != undefined)
    {
        var trace2 = {
            x: [plot_selected_x],
            y: [plot_selected_y],
            // mode: 'lines+markers',
            mode: 'markers',
            name: 'Data Point',
            line: {shape: 'linear'},
            type: 'scatter'
        };

        data = [trace1, trace2];
    }

    var layout = {
        title: title,
        legend: {
            y: 0.5,
            //traceorder: 'reversed',
            font: {size: 16},
            yref: 'paper'
        }};

    Plotly.newPlot(element_id, data, layout);
}