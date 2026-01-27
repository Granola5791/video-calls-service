import styled from "@emotion/styled";
import { Button } from "@mui/material";
import { Colors } from "../constants/general-contants.ts";


export const LongButton = styled(Button)({
    width: "300px",
    height: "50px",
    fontSize: "20px",
    fontWeight: "bold",
    borderRadius: "20px",
    border: `2px solid lightgray`,
    color: "gray",
});

export const LongButtonFilled = styled(LongButton)({
    color: "white",
    backgroundColor: Colors.primary,
    "&:hover": {
        backgroundColor: Colors.primaryHover,
    },
});

export const MediumButton = styled(Button)({
    width: "200px",
    height: "40px",
    fontSize: "16px",
    fontWeight: "bold",
    borderRadius: "20px",
    border: `2px solid lightgray`,
    color: "gray",
});

export const LeaveMeetingButton = styled(MediumButton)({
    color: "white",
    backgroundColor: Colors.danger,
    "&:hover": {
        backgroundColor: Colors.dangerHover,
    },
});