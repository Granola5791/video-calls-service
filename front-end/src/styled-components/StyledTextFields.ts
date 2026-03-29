import styled from "@emotion/styled";
import { TextField } from "@mui/material";

export const StyledTextField = styled(TextField)({
    width: "300px",
    textAlign: "center",
    '& .MuiInputLabel-root': {
        // Position the label on the right instead of the left
        right: '16px',
        left: 'auto',
        // Ensure it shrinks toward the right side
        transformOrigin: 'top right',
    },
    '& .MuiFilledInput-root': {
        // Fixes gray gaps in the background of the text field
        backgroundColor: '#e3f2fd',
        '&.Mui-focused': {
            backgroundColor: '#e3f2fd',
        },
    },
});