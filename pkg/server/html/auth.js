class Auth {
    constructor() {
        this.handlers = {};
        this.token = null;
        this.tokenKey = "webrddToken";
        this.authUrl = "/v1/login";
    }

    async login(message) {
        const token = localStorage.getItem(this.tokenKey);
        if (token) {
            console.log("Using previously stored token:", token)
            this.token = token;
            this.emit("login", {token: this.token});
            return true;
        }
        let userpass = await this.getUserPass(message);
        let username = userpass.user;
        let password = userpass.pass;
        let response;
        try {
            response = await fetch (this.authUrl, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({ username, password })
            });
        } catch (e) {
            console.log("Login attempt exception:", e);
            alert(e);
            this.emit("failure", {reason: e});
            return false;
        }
        console.log("Login status:", response.status, response.statusText);
        if (response.status > 299) {
            this.emit("failure", {reason: `${response.status} ${response.statusText}`});
            return false;
        }
        let result = await response.json();
        this.token = result.token;
        console.log("Got token:", this.token);
        let validFor =
            JSON.parse(atob(this.token.split(".")[1])).exp - Date.now() / 1000;
        console.log("Token expires in", validFor, "seconds");
        localStorage.setItem(this.tokenKey, this.token);
        this.emit("login", {token: this.token});
        return true;
    }

    async getUserPass(message) {
        let curtain = document.createElement("div");
        curtain.style.position = "fixed";
        curtain.style.top = "0px";
        curtain.style.left = "0px";
        curtain.style.width = "100%";
        curtain.style.height = "100%";
        curtain.style.backgroundColor = "#000000a0";
        curtain.style.display = "flex";
        curtain.style.alignItems = "center";
        curtain.style.justifyContent = "center";
    
        let login = document.createElement("form");
        login.style.backgroundColor = "#202020";
        login.style.color = "white";
        login.style.width = "250px";
        login.style.height = "100px";
        login.style.border = "5px solid";
        login.style.borderRadius = "20px";
        login.style.padding = "20px";
        login.style.display = "grid";
        login.style.alignItems = "center";
        login.style.justifyContent = "center";
        login.style.gridTemplateColumns = "1fr 2fr";

        if (message) {
            let messageDiv = document.createElement("div");
            messageDiv.innerHTML = message;
            messageDiv.style.gridColumn = "span 2";
            messageDiv.style.textAlign="center";
            messageDiv.style.fontWeight="bold";
            login.appendChild(messageDiv);
        }

        let usernameLabel = document.createElement("div");
        usernameLabel.innerHTML = "Username:";
        usernameLabel.style.padding = "0px 10px 0px 0px";
        login.appendChild(usernameLabel);
    
        let username = document.createElement("input");
        username.style.width = "100%";
        username.style.border = "1px solid white";
        username.style.backgroundColor = "black";
        username.style.color = "white";
        username.type = "text";
        login.appendChild(username);
    
        let passwordLabel = document.createElement("div");
        passwordLabel.innerHTML = "Password:";
        passwordLabel.style.padding = "0px 10px 0px 0px";
        login.appendChild(passwordLabel);
    
        let password = document.createElement("input");
        password.style.width = "100%";
        password.style.border = "1px solid white";
        password.style.backgroundColor = "black";
        password.style.color = "white";
        password.type = "password";
        login.appendChild(password);
    
        login.appendChild(document.createElement("div"));
    
        let submitContainer = document.createElement("div");
        submitContainer.style.display = "flex";
        submitContainer.style.justifyContent = "flex-end";
        let submitButton = document.createElement("input");
        submitButton.type = "submit";
        submitButton.style.border = "1px solid white";
        submitButton.style.backgroundColor = "black";
        submitButton.style.background =
          "radial-gradient(circle at center, blue 0px, black 100%)";
        submitButton.style.color = "white";
        submitButton.style.padding = "2px 10px";
        submitButton.value = "Login";
        submitContainer.appendChild(submitButton);
        login.appendChild(submitContainer);
    
        curtain.appendChild(login);
        document.body.appendChild(curtain);
    
        await new Promise((accept, reject) => {
            submitButton.onclick = () => accept();
        });
    
        curtain.remove();
    
        let user = username.value;
        let pass = password.value;
    
        return { user, pass };
    }

    reset() {
        this.token = null;
        localStorage.removeItem(this.tokenKey);
    }

    // Event handling
    on(eventName, handler) {
        if (!this.handlers[eventName]) {
            this.handlers[eventName] = [];
        }
        this.handlers[eventName].push(handler);
    }

    async emit(eventName, data) {
        // This waits for the next event loop so that we don't end up
        // with a really deep stack if any of the handlers call back
        // into a function that itself can call emit().
        await new Promise(r => {setTimeout(()=>r(), 0)});

        if (!this.handlers[eventName]) {
            this.handlers[eventName] = [];
        }
        for (const handler of this.handlers[eventName]) {
            handler(data);
        }
    }
    
}