import styled from "@emotion/styled";
import { Colors } from "../constants/general-contants.ts";

export const StyledTitle = styled.h1 ({
    fontSize: '3rem',
    fontWeight: 'bold',
    color: Colors.primary,
    textAlign: 'center',
});

export const BoldText = styled.p ({
    fontWeight: 'bold',
});

export const StyledH1 = styled.h1 ({
    color: Colors.primary,
});

export const MeetingID = styled.p ({
    position: 'absolute',
    right: '10px',
    color: 'gray',
    width: '30%',
    fontSize: '0.9rem',
});