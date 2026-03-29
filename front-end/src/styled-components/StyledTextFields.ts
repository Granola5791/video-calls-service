import styled from "@emotion/styled";
import { TextField } from "@mui/material";

export const StyledTextField = styled(TextField)({
    width: "300px",
    textAlign: "center",

    '& .MuiFilledInput-root': {
        // Fixes gray gaps in the background of the text field
        backgroundColor: '#e3f2fd',
        '&.Mui-focused': {
            backgroundColor: '#e3f2fd',
        },
        '&:hover': {
            backgroundColor: '#e3f2fd',
        },
    },
});