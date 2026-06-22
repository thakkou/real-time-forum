export const Header = () => (`
    <header class="navbar">
        <div class="logo">
            <a onclick="navigate('/')" style="text-decoration: none; color: inherit;">01Forum</a>
        </div>
        <div class="auth-buttons">
            {{if .IsLoggedIn}}
            <span class="welcome">Welcome, {{.User.Name}}</span>

            <button class="btn chat" data-count="9" onclick="navigate('chat')">
                <i class="fa-regular fa-message"></i>
            </button>

            <form action="/logout" method="POST">
                <button type="submit" class="btn logout">Logout</button>
            </form>
            {{else}}
            <a type="button" onclick="navigate('/login')" class="btn login">Login</a>
            <a type="button" onclick="navigate('/register')" class="btn register">Register</a>
            {{end}}
        </div>
    </header>
`);