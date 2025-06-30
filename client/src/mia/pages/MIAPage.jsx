import { useContext } from "react";
import { MIALayout } from "../layout/MIALayout";
import { Card, CardActions, CardContent, Container, Grid, Toolbar, Typography } from "@mui/material";
import { AppContext } from "../../context/AppContext";
import { DiskCard } from "../components/DiskCard";
import { Outlet } from "react-router-dom";

export const MIAPage = ({children}) => {

    const { disks } = useContext(AppContext);

    return (
        // <MIALayout>
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
                    { disks && disks.map((disk, index) => (
                        <DiskCard disk={disk} key={index} />
                    ))}
                </Grid>
            </Container>
            // <Outlet/>
        // </MIALayout>
    )
}
