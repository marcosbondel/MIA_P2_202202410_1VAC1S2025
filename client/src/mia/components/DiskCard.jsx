import { Button, Card, CardActions, CardContent, Grid, Typography } from "@mui/material"
import { useContext } from "react";
import { useNavigate } from "react-router-dom";
import { AppContext } from "../../context/AppContext";

export const DiskCard = ({disk}) => {
    const navigate = useNavigate();

    const handleDiskClick = () => {
        navigate(`/mia/disks/${disk.name}`);
    }

    return (
        <Grid xs={12} sm={6} md={4} lg={3}>
            <Card sx={{ minWidth: 275 }}>
                <CardContent>
                    <Typography variant="div">Disk</Typography>
                    <Typography variant="h3" align="center" sx={{ marginTop: 2, marginBottom: 2 }}>{ disk.name }</Typography>
                    <Typography variant="p">
                        <strong>Size:</strong> {disk.size}
                        <br />
                        <strong>Signature:</strong> {disk.signature}
                        <br />
                        <strong>Fit:</strong> {disk.fit}
                        <br />
                        <strong>Partitions:</strong> {disk.partitions?.length}
                    </Typography>
                </CardContent>
                <CardActions>
                    <Button 
                        variant="outlined"
                        fullWidth
                        onClick={handleDiskClick}
                    >
                        Partitions
                    </Button>
                </CardActions>
            </Card>
        </Grid>
    )
}
