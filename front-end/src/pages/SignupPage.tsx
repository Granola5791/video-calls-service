import React from 'react'
import { General, Auth, Errors, SuccessMessages } from "../constants/hebrew-constants"
import { CenteredColumn, CenteredFilledScreen } from "../styled-components/StyledBoxes"
import { LongButtonFilled } from "../styled-components/StyledButtons"
import { StyledTitle } from "../styled-components/StyledText"
import { StyledTextField } from "../styled-components/StyledTextFields"
import { BackendAddress, ApiEndpoints, HttpStatusCodes } from '../constants/backend-constants'
import { ErrorText } from '../styled-components/StyledErrors'
import { Typography } from '@mui/material'
import { AuthRequirements } from '../constants/general-contants'

const SignupPage = () => {

    const [username, setUsername] = React.useState<string>('');
    const [password, setPassword] = React.useState<string>('');
    const [rePassword, setRePassword] = React.useState<string>('');
    const [response, setResponse] = React.useState<string>('');

    const PasswordsMatch = (pass: string, rePass: string): boolean => {
        return pass === rePass;
    }

    const IsPasswordValid = (pass: string): boolean => {
        return pass.length >= AuthRequirements.passwordMinLength &&
            pass.length <= AuthRequirements.passwordMaxLength;
    }

    const IsUsernameValid = (username: string): boolean => {
        return username.length >= AuthRequirements.usernameMinLength &&
            username.length <= AuthRequirements.usernameMaxLength;
    }

    const HandleLogin = async () => {
        setResponse(Auth.wait);
        if (!IsUsernameValid(username)) {
            setResponse(Errors.invalidUsername);
            return;
        }
        if (!PasswordsMatch(password, rePassword)) {
            setResponse(Errors.passwordsDoNotMatch);
            return;
        }
        if (!IsPasswordValid(password)) {
            setResponse(Errors.invalidPasswordFormat);
            return;
        }

        const res = await fetch(BackendAddress + ApiEndpoints.signUp, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ username, password })
        })
        switch (res.status) {
            case HttpStatusCodes.Created:
                setResponse(SuccessMessages.userCreated);
                break;
            case HttpStatusCodes.Conflict:
                setResponse(Errors.usernameAlreadyExists);
                break;
            case HttpStatusCodes.BadRequest:
                setResponse(Errors.invalidPasswordFormat);
                break;
            default:
                setResponse(Errors.genericError);
        }
    }

    return (
        <CenteredFilledScreen>
            <CenteredColumn >
                <StyledTitle>{General.appName}</StyledTitle>

                <CenteredColumn>
                    <Typography>{Auth.passwordRequirements}</Typography>
                    <StyledTextField
                        label={Auth.username}
                        variant="filled"
                        onChange={(e) => setUsername(e.target.value)}
                    />
                    <StyledTextField
                        label={Auth.password}
                        type="password"
                        variant="filled"
                        onChange={(e) => setPassword(e.target.value)}
                    />
                    <StyledTextField
                        label={Auth.rePassword}
                        type="password"
                        variant="filled"
                        onChange={(e) => setRePassword(e.target.value)}
                    />
                    <LongButtonFilled onClick={HandleLogin}>{Auth.signUpButton}</LongButtonFilled>
                    <ErrorText>{response}</ErrorText>
                </ CenteredColumn>
            </CenteredColumn>
        </CenteredFilledScreen>
    )
}

export default SignupPage