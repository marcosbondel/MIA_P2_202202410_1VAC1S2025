import { Box, Typography } from "@mui/material";
import { NavBar } from "../components/NavBar";

const drawerWidth = 240

export const MIALayout = ({children}) => {
    return (
        <Box
            sx={{ display: 'flex' }}
        >
            {/* <Sidebar /> */}
            {/* <Sidebar drawerWidth={drawerWidth}/> */}
            <NavBar drawerWidth={drawerWidth} />
            {/* <NavBar drawerWidth={drawerWidth}/> */}
            {/* <NavBar /> */}
            {children}
        </Box>
    )
}
