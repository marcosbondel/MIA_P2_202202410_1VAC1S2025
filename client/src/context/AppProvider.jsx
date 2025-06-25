import { useEffect, useReducer } from "react";
import { AppContext } from "./AppContext"
import { appReducer } from "./appReducer";
import { Alert, Snackbar } from "@mui/material";
import { useNavigate } from "react-router-dom";
import { useLocation } from "react-router-dom";
import queryString from "query-string"

const init = () => {
    return {
        logged: false,
        error: "",
        successMessage: "",
        result: {},
        showError: false,
        showSuccess: false,
        disks: [],
        partitions: [],
        current_fs_location: '/',
        current_fs: '',
        current_directory: '',
        directory_parts: '',
        current_file_content: '',
    }
}

export const AppProvider = ({ children }) => {
    const [state, dispatch] = useReducer(appReducer, {}, init)
    const navigate = useNavigate();

    const login = async(id, username, password) => {
        const response = await fetch('http://localhost:3000/api/auth/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ Id: id, User: username, Pass: password })
        })

        if(!response.ok) {
            const error = await response.json();
            dispatch({ type: 'error[set]', payload: { error: "Login failed :("} });
            return;
        }

        const data = await response.json();
        console.log("Login successful:", data);
        dispatch({ type: 'logged[set]', payload: { logged: true } });

        navigate('/mia')
    }

    const getDisks = async() => {
        const response = await fetch('http://localhost:3000/api/disks');
        if(!response.ok) {
            const error = await response.json();
            console.error("Failed to fetch disks:", error);
            return
        }

        const disks = await response.json();

        dispatch({ type: 'disks[set]', payload: disks });
    }

    const getPartitions = async(disk) => {
        const response = await fetch(`http://localhost:3000/api/disks/${disk}/partitions`);
        if(!response.ok) {
            const error = await response.json();
            console.error("Failed to fetch partitions:", error);
            dispatch({ type: 'error[set]', payload: { error: "Failed to fetch partitions" } });
            return
        }

        const partitions = await response.json();
        dispatch({ type: 'partitions[set]', payload: {partitions} });
    }

    const logout = async() => {
        const response = await fetch('http://localhost:3000/api/auth/logout', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({})
        })

        if(!response.ok){
            dispatch({ type: 'error[set]', payload: { error: "Logout failed :("} });
            return
        }

        dispatch({ type: 'logged[set]', payload: { logged: false } });

        navigate('/login');
    }

    const getFileSystem = async() => {
        let url = 'http://localhost:3000/api/fs'

        if(state.current_fs_location != '/'){
            url = `${url}?path=${state.current_fs_location}`;
        }else{
            url = `${url}?path=/`;
        }

        const response = await fetch(url)

        if(!response.ok) {
            const error = await response.json();
            console.error("Failed to fetch filesystem:", error);
            dispatch({ type: 'error[set]', payload: { error: "Failed to fetch filesystem" } });
            return
        }
        const data = await response.json();
        console.log(data)
        dispatch({ type: 'current_fs[set]', payload: { current_fs: data } });
        dispatch({ type: 'current_directory[set]', payload: { current_directory: data?.children } });
        dispatch({ type: 'directory_parts[set]', payload: { directory_parts: data.path.split("/").slice(1) } });
        if(data?.type == 'file'){
            dispatch({ type: 'current_file_content[set]', payload: { current_file_content: data.content } });
        }else{
            dispatch({ type: 'current_file_content[set]', payload: { current_file_content: '' } });
        }
    }
    
    const dispatchCurrentFSLocation = (path) => {
        dispatch({ type: 'current_fs_location[set]', payload: { current_fs_location: path } });
    }

    const showSuccessMessage = (message) => {
        dispatch({ type: 'showSuccess[set]', payload: { successMessage: message } });
    }

    useEffect(() => {
        getDisks()
    }, [])
    

    return (
        <AppContext.Provider value={{
            ...state,
            login,
            getDisks,
            logout,
            getPartitions,
            getFileSystem,
            dispatchCurrentFSLocation,
            showSuccessMessage
        }}>
            {children}
            <Snackbar
                open={state.showError}
                autoHideDuration={3000}
                onClose={() => dispatch({ type: 'error[set]', payload: { error: "" } })}
                message={state.error}
                anchorOrigin={{ vertical: 'top', horizontal: 'right' }}
                color="error"
            />
            <Snackbar
                open={state.showSuccess}
                autoHideDuration={3000}
                onClose={() => dispatch({ type: 'showSuccess[set]', payload: { successMessage: "" } })}
                message={state.successMessage}
                anchorOrigin={{ vertical: 'top', horizontal: 'right' }}
                color="info"
            />
        </AppContext.Provider>
    )
}
