import { Card, CardActions, CardContent, Grid, Typography } from "@mui/material"
import { useNavigate } from "react-router-dom";

export const DiskCard = ({disk}) => {
    const navigate = useNavigate();

    const handleDiskClick = () => {
        navigate(`/mia/disks/${disk}`);
    }

    return (
        <Grid item key={disk} xs={12} sm={6} md={4} lg={3} onClick={handleDiskClick}>
            <Card sx={{ minWidth: 275 }}>
                <CardContent>
                    <Typography>Disk</Typography>
                    <Typography variant="h5" align="center">{ disk }</Typography>
                </CardContent>
                <CardActions>
                    <Typography variant="body2">Size: 0</Typography>
                </CardActions>
            </Card>
        </Grid>
    )
}
