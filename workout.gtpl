<html>
    <head>
        <title>Workout</title>
    </head>
    <body>
        <form action="/workout" method="post">
            Reps: <input type="text" name="reps">
            Weight: <input type="text" name="weight">
            <input type="hidden" value="{{.}}">
            <input type="submit" value="Login">
        </form>
    </body>
</html>

