import { Box, Toolbar, Typography } from "@mui/material";
import { NavBar } from "../components/NavBar";
import { Outlet } from "react-router-dom";

const drawerWidth = 240

export const MIALayout = ({children}) => {
    return (
        <Box
            sx={{ display: 'flex' }}
        >
            <NavBar drawerWidth={drawerWidth} />
            
            <Outlet/>
        </Box>
    )
}
