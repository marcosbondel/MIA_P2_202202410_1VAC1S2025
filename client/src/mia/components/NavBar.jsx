
import { MenuOutlined, Terminal } from "@mui/icons-material"
import { AppBar, Box, Button, Grid, IconButton, Link, Toolbar, Typography } from "@mui/material"
import { useContext } from "react"
import { AppContext } from "../../context/AppContext"
import { useNavigate, Link as RouterLink } from "react-router-dom"

export const NavBar = () => {
    const navigate = useNavigate()

    const { logout } = useContext(AppContext)

    const onLogout = async (e) => {
        e.preventDefault()
        await logout()
    }

    const goToTerminal = () => {
        console.log("/mia/terminal")
        navigate('/mia/terminal')
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
                    <Link component={RouterLink} color='inherit' to='/mia' underline="none">
                        <Typography variant="h6" noWrap component='div'>MIA</Typography>
                    </Link>
                    <Grid 
                        item
                        sx={{
                            display: 'flex',
                            alignItems: 'center',
                            justifyContent: 'center',
                            gap: 2,
                        }}
                    >
                        <IconButton onClick={goToTerminal}>
                            <Terminal sx={{ fontSize: '3rem', color: 'primary.main' }}/>
                        </IconButton>
                        <form onSubmit={onLogout}>
                            <Button variant="contained" type="submit">Logout</Button>
                        </form>
                    </Grid>
                </Grid>
            </Toolbar>
        </AppBar>
    )
}
