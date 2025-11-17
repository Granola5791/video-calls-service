import React from 'react'
import { General, Auth, Errors, SuccessMessages } from "../constants/hebrew-constants"
import { CenteredColumn } from "../styled-components/StyledBoxes"
import { LongButtonFilled } from "../styled-components/StyledButtons"
import { StyledTitle } from "../styled-components/StyledText"
import { StyledTextField } from "../styled-components/StyledTextFields"
import { BackendAddress, ApiEndpoints, HttpStatusCodes } from '../constants/backend-constants'
import { ErrorText } from '../styled-components/StyledErrors'

const SignupPage = () => {

    const [username, setUsername] = React.useState<string>('');
    const [password, setPassword] = React.useState<string>('');
    const [rePassword, setRePassword] = React.useState<string>('');
    const [response, setResponse] = React.useState<string>('');

    const PasswordsMatch = (pass: string, rePass: string): boolean => {
        return pass === rePass;
    }

    const HandleLogin = async () => {
        if (!PasswordsMatch(password, rePassword)) {
            setResponse(Errors.passwordsDoNotMatch);
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
            default:
                setResponse(Errors.genericError);
        }
    }

    return (
        <>
            <StyledTitle>{General.appName}</StyledTitle>

            <CenteredColumn>
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
                <LongButtonFilled onClick={HandleLogin}>{Auth.loginButton}</LongButtonFilled>
                <ErrorText>{response}</ErrorText>
            </ CenteredColumn>

        </>
    )
}

export default SignupPage