export const Header = (nickname = "NONE") => (`
    <header class="navbar">
        <div class="logo">
            <a onclick="navigate('/')" style="text-decoration: none; color: inherit;">01Forum</a>
        </div>
        <div class="auth-buttons">
            <span class="welcome">Welcome, ${nickname}</span>

            <button class="btn chat" data-count="9" onclick="navigate('/chat')">
                <i class="fa-regular fa-message"></i>
            </button>
            
            <button class="btn logout" id="logout-btn">
                Logout
            </button>
        </div>
    </header>
`);