import { useContext } from "react"
import { AppContext } from "../context/AppContext"
import { Navigate, Outlet } from "react-router-dom"

export const PrivateRoute = ({children}) => {
    const { logged } = useContext(AppContext)

    return (logged)
    ? children
    : <Navigate to="/login" />
}
