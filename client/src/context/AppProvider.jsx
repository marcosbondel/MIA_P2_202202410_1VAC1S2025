import { useEffect, useReducer } from "react";
import { AppContext } from "./AppContext"
import { appReducer } from "./appReducer";
import { Alert, Snackbar } from "@mui/material";
import { useNavigate } from "react-router-dom";

const init = () => {
    return {
        error: "",
        result: {},
        showError: false,
        disks: JSON.parse(localStorage.getItem('disks')) || []
    }
}

export const AppProvider = ({ children }) => {
    const [state, dispatch] = useReducer(appReducer, {}, init)
    const navigate = useNavigate();

    const onLogin = async(id, username, password) => {
        console.log(username, password);
        dispatch({ type: 'error[set]', payload: { error: "" } });
        const response = await fetch('http://localhost:3000/api/auth/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ Id: id, User: username, Pass: password })
        })

        if(!response.ok) {
            const error = await response.json();
            console.error("Login failed:", error);
            dispatch({ type: 'error[set]', payload: { error: "Login failed :("} });
            return;
        }

        const data = await response.json();
        console.log("Login successful:", data);
        navigate('/mia')
    }

    const getDisks = async() => {
        const response = await fetch('http://localhost:3000/api/disks');
        if(!response.ok) {
            const error = await response.json();
            console.error("Failed to fetch disks:", error);
            return [];
        }

        const disks = await response.json();
        console.log("Disks fetched successfully:", disks);
        localStorage.setItem('disks', JSON.stringify(disks));

        dispatch({ type: 'disks[set]', payload: { disks } });
        return disks;
    }

    useEffect(() => {

        if (Object.keys(state.disks).length !== 0) return

        getDisks()

    }, [])
    

    return (
        <AppContext.Provider value={{
            ...state,
            onLogin,
            getDisks,
        }}>
            {children}
            {/* { state.error.length != 0 && */}
                {/* <Alert severity="warning">This is a warning Alert.</Alert> */}
            {/* } */}
            <pre>{ state.showError }</pre>
            <Snackbar
                open={state.showError}
                autoHideDuration={3000}
                onClose={() => dispatch({ type: 'error[set]', payload: { error: "" } })}
                message={state.error}
                anchorOrigin={{ vertical: 'top', horizontal: 'right' }}
                color="error"
            />
        </AppContext.Provider>
    )
}
