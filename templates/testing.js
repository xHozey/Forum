if (auth) {
    document.write(`
        <div id="session">
            <a href="/logout"><button>Log Out</button></a>
            <h1>Welcome, ${username}!</h1>
            <h3>Create a new post</h3>
            <form method="post" action="/post">
                <textarea name="textarea" rows="5" cols="40" placeholder="Write something..."></textarea>
                <br>
                <button type="submit">Post</button>
            </form>
        </div>
    `);
} else {
    document.write(`<div id="guest">
        <a href="/login"><button>log in</button></a>
        <a href="/register"><button>register</button></a>
    </div>`)
}