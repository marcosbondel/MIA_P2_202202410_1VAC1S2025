import { Button, Card, CardActions, CardContent, Grid, Typography } from "@mui/material"
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
                    <Typography variant="h5" align="center" sx={{ marginTop: 2, marginBottom: 2 }}>{ partition.name }</Typography>
                    <Typography variant="p">
                        <strong>Mounted:</strong> {partition.mounted}
                        <br />
                        <strong>Size:</strong> {partition.size}
                        <br />
                        <strong>Fit:</strong> {partition.fit}
                        <br />
                        <strong>Type:</strong> {partition.type}
                    </Typography>
                </CardContent>
                <CardActions>
                    <Button 
                        variant="outlined"
                        fullWidth
                        onClick={handlePartitionClick}
                    >
                        FS
                    </Button>
                </CardActions>
            </Card>
        </Grid>
    )
}
