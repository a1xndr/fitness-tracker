<html>
    <head>
        <title>Workout</title>
    </head>
    <body>
    <h1> Workout on {{.FormattedDate}}</h1>
    <table>
    {{range .Sets}}
        <tr>
            <td>{{.Exercise}}</td>
            <td>{{.Reps}}</td>
            <td>{{.Weight}}</td>
        </tr>
    {{end}}
    <pre>{{.FormatAsMd}}</pre>
    <pre>{{.FormatAsAsciiTable}}</pre>
    </table>
        <form action="/workout" method="post">
            Reps: <input type="text" name="reps">
            Weight: <input type="text" name="weight">
            <input type="hidden" value="{{.}}">
            <input type="submit" value="Login">
        </form>
    </body>
</html>

