import { AppContext } from "./AppContext"

export const AppProvider = ({ children }) => {

    const onLogin = async(e, username, password) => {
        const response = await fetch('http://localhost:3000/api/auth/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ username, password })
        })

        if(!response.ok) {
            const error = await response.json();
            console.error("Login failed:", error);
            return;
        }

        const data = await response.json();
        console.log("Login successful:", data);
    }

    return (
        <AppContext.Provider value={{
            onLogin
        }}>
            {children}
        </AppContext.Provider>
    )
}
