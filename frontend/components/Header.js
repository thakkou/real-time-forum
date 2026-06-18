export const Header = () => (`
    <header class="navbar">
        <div class="logo">
            <a href="/" style="text-decoration: none; color: inherit;">01Forum</a>
        </div>
        <div class="auth-buttons">
            {{if .IsLoggedIn}}
            <span class="welcome">Welcome, {{.User.Name}}</span>
            <form action="/logout" method="POST">
                <button type="submit" class="btn logout">Logout</button>
            </form>
            {{else}}
            <a type="button" href="/login" class="btn login">Login</a>
            <a type="button" href="/register" class="btn register">Register</a>
            {{end}}
        </div>
    </header>
`);