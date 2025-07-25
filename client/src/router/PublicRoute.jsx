import { useContext } from "react"
import { AppContext } from "../context/AppContext"
import { Navigate } from "react-router-dom"

export const PublicRoute = ({children}) => {
    const { logged } = useContext(AppContext)
    
    return (!logged) 
        ? children
        : <Navigate to="/mia"/>
}
