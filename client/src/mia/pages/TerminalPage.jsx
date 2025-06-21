import { Container, Toolbar, Box, TextField } from "@mui/material"
import { useState, useRef, useEffect } from "react"
import { MIALayout } from "../layout/MIALayout"

export const TerminalPage = () => {
    const [lines, setLines] = useState(["Welcome to MIA Terminal. Type a command."]);
    const [input, setInput] = useState("");
    const terminalEndRef = useRef(null);

    const runCommand = async(command) => {
        // Placeholder for command execution logic
        console.log(`Executing command: ${command}`);
        const response = await fetch(`http://localhost:3000/api/run_command`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ command_string: command }),
        });
        if (!response.ok) {
            console.log(`Error executing command: ${command}`);
            return;
        }

        const data = await response.json();
        setLines((prev) => [...prev, `$ ${command}`, data.output]);
        console.log(`Command executed: ${command}`);
        console.log(`Command output: ${data}`);
    }

    const handleCommand = async(cmd) => {
        // Example command handling
        let output;
        switch (cmd.trim()) {
            case "help":
                output = "Available commands: help, clear, echo";
                break;
            case "clear":
                setLines([]);
                return;
            case "":
                output = "";
                break;
            default:
                await runCommand(cmd);
                // output = `Command not found: ${cmd}`;
        }

        setLines((prev) => [...prev, `$ ${cmd}`, output]);
    };

    const handleKeyDown = (e) => {
        if (e.key === "Enter") {
            e.preventDefault();
            handleCommand(input);
            setInput("");
        }
    };

    useEffect(() => {
        terminalEndRef.current?.scrollIntoView({ behavior: "smooth" });
    }, [lines]);

    return (
        <MIALayout>
            <Container >
                <Toolbar />
                <Toolbar />
                <Box
                    sx={{
                        backgroundColor: "#1e1e1e",
                        color: "#d4d4d4",
                        padding: 2,
                        borderRadius: 2,
                        height: "70vh",
                        overflowY: "auto",
                        fontFamily: "monospace",
                        whiteSpace: "pre-wrap",
                        display: "flex",
                        flexDirection: "column",
                    }}
                >
                    {lines.map((line, idx) => (
                        <Box key={idx}>{line}</Box>
                    ))}

                    <Box sx={{ display: "flex" }}>
                        <Box sx={{ mr: 1 }}>$</Box>
                        <TextField
                            variant="standard"
                            fullWidth
                            value={input}
                            onChange={(e) => setInput(e.target.value)}
                            onKeyDown={handleKeyDown}
                            InputProps={{
                                disableUnderline: true,
                                style: { color: "#d4d4d4", fontFamily: "monospace" },
                            }}
                        />
                    </Box>
                    <div ref={terminalEndRef} />
                </Box>
            </Container>
        </MIALayout>
    );
}
