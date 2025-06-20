
import { LogoutOutlined, MenuOutlined } from "@mui/icons-material"
import { AppBar, Box, Button, Grid, IconButton, Toolbar, Typography } from "@mui/material"
import { useContext } from "react"
import { AppContext } from "../../context/AppContext"

export const NavBar = () => {

    const { logout } = useContext(AppContext)

    const onLogout = async (e) => {
        e.preventDefault()
        await logout()
    }

    return (
        <AppBar 
            position="fixed"
            sx={{
                backgroundColor: 'secondary.main'
            }}
        >
            <Toolbar>
                <IconButton
                    color='inherit'
                    edge='start'
                    sx={{ mr: 2, display: { sm: 'none' } }}
                >
                    <MenuOutlined/>
                </IconButton>
                <Grid
                    container
                    direction="row"
                    sx={{
                        flexGrow: 1,
                        justifyContent: "space-between",
                        alignItems: "center",
                    }}
                >
                    <Typography variant="h6" noWrap component='div'>MIA</Typography>
                    {/* <IconButton color='info' size="large">
                        <LogoutOutlined size="large"/>
                    </IconButton> */}
                    <form onSubmit={onLogout}>
                        <Button variant="contained" type="submit">Logout</Button>
                    </form>
                </Grid>
            </Toolbar>
        </AppBar>
    )
}
