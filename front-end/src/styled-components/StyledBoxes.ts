import styled from "@emotion/styled";
import { Colors } from "../constants/general-contants";


export const CenteredColumn = styled.div({
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
    margin: 'auto',
    flexDirection: 'column',
    gap: '10px',
    backgroundColor: 'white',
    padding: '20px',
    borderRadius: '10px',
})

export const CenteredRow = styled.div({
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
    margin: 'auto',
    flexDirection: 'row',
    gap: '10px',
    backgroundColor: 'white',
    padding: '20px',
    borderRadius: '10px',
})

export const CenteredFilledScreen = styled.div({
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
    height: '100vh',
    width: '100vw',
    margin: 'auto',
    backgroundColor: Colors.primary,
})

export const StyledMeetingGrid = styled.div({
    display: "grid",
    position: "fixed",
    gridTemplateColumns: "repeat(auto-fit, minmax(250px, 1fr))",
    gap: "8px",
    width: "100%",
    height: "100%",
    top: "0px",
    right: "0px",
    left: "0px",
    bottom: "70px",
    padding: "8px",
    boxSizing: "border-box",
    backgroundColor: Colors.meetingBackground,
});