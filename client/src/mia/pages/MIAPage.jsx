import { MIALayout } from "../layout/MIALayout";
import { Card, CardActions, CardContent, Container, Grid, Toolbar, Typography } from "@mui/material";
// import { MIARoutes } from "../routes/MIARoutes";

export const MIAPage = ({children}) => {

    const disks = [
        { name: "C", size: "500GB" },
        { name: "D", size: "1TB" },
        { name: "E", size: "2TB" },
        { name: "F", size: "250GB" },
        { name: "G", size: "1.5TB" },
        { name: "H", size: "3TB" },
        { name: "I", size: "4TB" },
        { name: "J", size: "500GB" }
    ]

    return (
        <MIALayout>
            {/* <MIARoutes/> */}
            <Container>
                <Toolbar/>
                <h1>My Disks</h1>
                { children }
                <Grid
                    container
                    spacing={2}
                    direction="row"
                    justifyContent="space-between"
                    alignItems="center"
                >
                    { disks.map((disk) => (
                        <Grid key={disk.name} item xs={12} sm={6} md={4} lg={3}>
                            <Card sx={{ minWidth: 275 }}>
                                <CardContent>
                                    <Typography>Disk</Typography>
                                    <Typography variant="h5" align="center">{ disk.name }</Typography>
                                </CardContent>
                                <CardActions>
                                    <Typography variant="body2">Size: { disk.size }</Typography>
                                </CardActions>
                            </Card>
                        </Grid>
                    ))}
                </Grid>
            </Container>
        </MIALayout>
    )
}
