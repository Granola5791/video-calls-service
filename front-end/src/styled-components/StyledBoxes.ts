import styled from "@emotion/styled";
import { Colors } from "../constants/general-contants";


export const CenteredColumn = styled.div ({
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
    // height: '100%',
    margin: 'auto',
    flexDirection: 'column',
    gap: '10px',
    backgroundColor: 'white',
    padding: '20px',
    borderRadius: '10px',
})

export const CenteredFilledScreen = styled.div ({
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
    height: '100vh',
    width: '100vw',
    margin: 'auto',
    backgroundColor: Colors.primary,
})