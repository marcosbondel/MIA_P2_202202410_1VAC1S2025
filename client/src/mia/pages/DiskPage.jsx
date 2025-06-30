import { Container, Grid, Toolbar, Typography } from "@mui/material"
import { useContext, useEffect } from "react"
import { useParams } from "react-router-dom"
import { AppContext } from "../../context/AppContext"
import { PartitionCard } from "../components/PartitionCard"

export const DiskPage = () => {
    const {disk} = useParams()
    const { getPartitions, partitions } = useContext(AppContext)
    
    useEffect(() => {
        getPartitions(disk);
    }, [])

    return (
        <Container>
            <Toolbar/>
            <h1>Disk: {disk}</h1>
            <Grid
                container
                spacing={2}
                direction="row"
                justifyContent="left"
                alignItems="center"
            >
                { (partitions && partitions.length > 0 ? (
                    partitions.map((partition, index) => (
                        <PartitionCard key={`${partition.name}-${index}`} partition={partition} />
                    ))
                ) : (
                    <Typography>No partitions found for disk {disk}</Typography>
                )) }
            </Grid>
        </Container>
    )
}
