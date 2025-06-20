
import { LogoutOutlined, MenuOutlined } from "@mui/icons-material"
import { AppBar, Box, Button, Grid, IconButton, Toolbar, Typography } from "@mui/material"

export const NavBar = () => {
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
                    <Button variant="contained">Logout</Button>
                </Grid>
            </Toolbar>
        </AppBar>
    )
}
