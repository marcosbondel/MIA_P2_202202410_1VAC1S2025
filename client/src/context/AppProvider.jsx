import { useEffect, useReducer } from "react";
import { AppContext } from "./AppContext"
import { appReducer } from "./appReducer";
import { Alert, Snackbar } from "@mui/material";
import { useNavigate } from "react-router-dom";
import { useLocation } from "react-router-dom";
import queryString from "query-string"

import { url } from "../api/url";

const init = () => {
    return {
        logged: JSON.parse(localStorage.getItem('logged')) || false,
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
        fs_id: '',
    }
}

export const AppProvider = ({ children }) => {
    const [state, dispatch] = useReducer(appReducer, {}, init)
    const navigate = useNavigate();

    const login = async(id, username, password) => {
        const response = await fetch(`${url.base}/api/auth/login`, {
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
        localStorage.setItem('logged', JSON.stringify(true));

        navigate('/mia')
    }

    const getDisks = async() => {
        const response = await fetch(`${url.base}/api/disks`);
        if(!response.ok) {
            const error = await response.json();
            console.error("Failed to fetch disks:", error);
            return
        }

        const disks = await response.json();
        dispatch({ type: 'disks[set]', payload: disks });
        console.log(`Disks fetched successfully:`, disks);
    }

    const getPartitions = async(disk) => {
        const response = await fetch(`${url.base}/api/disks/${disk}/partitions`);
        if(!response.ok) {
            const error = await response.json();
            console.error("Failed to fetch partitions:", error);
            dispatch({ type: 'error[set]', payload: { error: "Failed to fetch partitions" } });
            return
        }

        const partitions = await response.json();
        dispatch({ type: 'partitions[set]', payload: {partitions} });
        console.log(`Partitions: `, partitions);
    }

    const logout = async() => {
        const response = await fetch(`${url.base}/api/auth/logout`, {
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
        localStorage.removeItem('logged');

        navigate('/login');
    }

    const getFileSystem = async() => {
        let endpoint = `${url.base}/api/fs?id=${state.fs_id}`;

        if(state.current_fs_location != '/'){
            endpoint = `${endpoint}&path=${state.current_fs_location}`;
        }else{
            endpoint = `${endpoint}&path=/`;
        }

        console.log(`Fetching filesystem from: ${endpoint}`);

        const response = await fetch(endpoint)

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

    const modifyFSId = (id) => {
        console.log(`Modifying fs_id to: ${id}`);
        dispatch({ type: 'fs_id[set]', payload: { fs_id: id } });
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
            showSuccessMessage,
            modifyFSId
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
