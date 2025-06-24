import { Card, CardActions, CardContent, Grid, Typography } from "@mui/material"
import { useNavigate, useParams } from "react-router-dom";

export const PartitionCard = ({partition}) => {

    const navigate = useNavigate();
    const { disk } = useParams()

    const handlePartitionClick = () => {
        navigate(`/mia/disks/${disk}/partitions/${partition.name}`);
    }
    return (
        <Grid item key={partition.name} xs={12} sm={6} md={4} lg={3} onClick={handlePartitionClick}>
            <Card sx={{ minWidth: 275 }}>
                <CardContent>
                    <Typography>Partition</Typography>
                    <Typography variant="h5" align="center">{ partition.name }</Typography>
                </CardContent>
                <CardActions>
                    <Typography variant="body2">Size: {partition.size}</Typography>
                </CardActions>
            </Card>
        </Grid>
    )
}
