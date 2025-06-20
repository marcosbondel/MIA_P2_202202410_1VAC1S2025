import { Button, Grid, TextField, Typography } from "@mui/material"
import { useForm } from "../../hooks"
import { useContext } from "react"
import { useNavigate } from "react-router-dom"
import { AppContext } from "../../context/AppContext"

export const LoginPage = () => {

    const { onLogin } = useContext(AppContext)
    const navigate = useNavigate()

    const {username, password, onInputChange} = useForm({
        username: 'root',
        password: 'pass'
    })

    const onSubmit = (e) => {
        e.preventDefault()
        onLogin(username, password)
        navigate('/mia')
    }

    return (
        <Grid 
            container
            spacing={ 0 }
            direction="column"
            alignItems="center"
            justifyContent="center"
            sx={{ minHeight: '100vh', padding: 4, backgroundColor: 'secondary.main' }}
        >
            <Grid 
                item
                className="box-shadow"
                xs={ 3 }
                sx={{ 
                    width: { md: 450},
                    backgroundColor: 'white', 
                    padding: 3, 
                    borderRadius: 2 
                }}
            >
                <Typography variant="h3" sx={{mb: 1, textAlign: 'center'}}>MIA - 202202410</Typography>

                <form onSubmit={onSubmit}>
                    <Grid container>
                        <Grid item size={ 12 } sx={{mt: 2}}>
                            <TextField 
                                label="Username" 
                                type="text" 
                                placeholder="Write your username"
                                name="username"
                                value={username}
                                onChange={onInputChange}
                                fullWidth
                            />

                        </Grid>
                        <Grid item size={ 12 } sx={{mt: 2}}>
                            <TextField 
                                label="Password" 
                                type="password"
                                placeholder="Write your password"
                                name="password"
                                value={password}
                                onChange={onInputChange}
                                fullWidth
                            />

                        </Grid>
                        <Grid item size={12} sx={{mt: 2}}>
                            <Button type="submit" variant="contained" fullWidth>Login</Button>
                        </Grid>
                    </Grid>

                </form>

            </Grid>
        </Grid>
    )
}
