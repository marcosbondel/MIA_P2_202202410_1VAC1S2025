import { Button, Grid, TextField, Typography } from "@mui/material"

export const LoginPage = () => {
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

                <form>
                    <Grid container>
                        <Grid item size={ 12 } sx={{mt: 2}}>
                            <TextField 
                                label="Username" 
                                type="text" 
                                placeholder="root"
                                name="username"
                                value="username"
                                fullWidth
                            />

                        </Grid>
                        <Grid item size={ 12 } sx={{mt: 2}}>
                            <TextField 
                                label="Password" 
                                type="password"
                                placeholder="pass"
                                name="password"
                                value="password"
                                fullWidth
                            />

                        </Grid>
                        <Grid item size={12} sx={{mt: 2}}>
                            <Button variant="contained" fullWidth>Login</Button>
                        </Grid>
                    </Grid>

                </form>

            </Grid>
        </Grid>
    )
}
