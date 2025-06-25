import { Button, Grid, InputLabel, MenuItem, Select, TextField, Typography } from "@mui/material"
import { useForm } from "../../hooks"
import { useContext, useEffect } from "react"
import { useNavigate } from "react-router-dom"
import { AppContext } from "../../context/AppContext"
import { TerminalPage } from "../../mia/pages"
import { UploadFile } from "../../mia/components/Upload"

export const LoginPage = () => {

    const { login, disks } = useContext(AppContext)
    const navigate = useNavigate()

    const {id, username, password, onInputChange} = useForm({
        id: '',
        username: '',
        password: '',
    })

    const onSubmit = (e) => {
        e.preventDefault()
        login(id, username, password)
        // navigate('/mia')
    }

    return (
        <>
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
                                    label="ID" 
                                    type="text" 
                                    placeholder="ID"
                                    name="id"
                                    value={id}
                                    onChange={onInputChange}
                                    fullWidth
                                />
                            </Grid>
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
            <Grid
                item
                sx={{ 
                    backgroundColor: 'white', 
                    padding: 3, 
                    borderRadius: 2, 
                }}
            >
                <TerminalPage/>
                <br />
                <UploadFile/>
            </Grid>
        
        </>
    )
}
