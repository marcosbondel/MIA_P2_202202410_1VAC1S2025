import { Button, Grid, TextField } from "@mui/material";
import { useContext, useState } from "react";
import { AppContext } from "../../context/AppContext";

import { url } from "../../api/url";

export const UploadFile = () => {
    const [file, setFile] = useState(null);
    const { showSuccessMessage } = useContext(AppContext);

    const handleSubmit = async () => {
        if (!file) return;

        const formData = new FormData();
        formData.append("file", file);

        const res = await fetch(`${url.base}/api/upload`, {
        method: "POST",
        body: formData,
        });

        const data = await res.json();
        showSuccessMessage("SDAA file uploaded successfully!");
    };

    return (
        <Grid
            container
            spacing={2}
            direction="column"
            alignItems="center"
            justifyContent="center"
        >
            <TextField type="file" accept=".sdaa" onChange={(e) => setFile(e.target.files[0])}/>
            <Button onClick={handleSubmit} variant="contained">Upload SDAA</Button>
        </Grid>
    );
}
