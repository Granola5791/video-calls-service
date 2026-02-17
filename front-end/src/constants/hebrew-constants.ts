import { AuthRequirements } from "./general-contants";

export const General = {
    appName: "שו״צ (זמני)",
};

export const Auth = {
    wait: "רגע...",
    username: "שם משתמש",
    password: "סיסמה",
    rePassword: "אימות סיסמה",
    loginButton: "התחבר",
    signUpButton: "הרשמה",
    passwordRequirements: `הסיסמה חייבת להכיל בין ${AuthRequirements.passwordMinLength} ל־${AuthRequirements.passwordMaxLength} תווים.`,
    alreadyHaveAccount: {
        text: "כבר יש לכם חשבון? ",
        linkText: "להתחברות!",
    },
    noAccountYet: {
        text: "עוד אין לכם חשבון? ",
        linkText: "להרשמה!",
    },
    goToLogin: "להתחברות!",
};

export const Errors = {
    passwordsDoNotMatch: "הסיסמאות לא תואמות!",
    usernameAlreadyExists: "שם משתמש זה תפוס!",
    genericError: "אירעה שגיאה, אנא נסו שוב במועד מאוחר יותר",
    invalidPasswordFormat: "הסיסמה לא עומדת בדרישות!",
    invalidUsername: "יש להכניס שם משתמש תקין!",
    invalidAuthInput: "שם המשתמש או הסיסמה שגויים!",
};

export const SuccessMessages = {
    userCreated: "המשתמש נוצר בהצלחה",
    userLoggedIn: "המשתמש התחבר בהצלחה",
};

export const LandingPageText = {
    joinMeetingButton: "הצטרף לפגישה",
    signUpButton: "הרשמה",
    logInButton: "התחברות",
};

export const Stream = {
    startButton: "התחל שידור",
    stopButton: "עצור שידור",
}

export const HomePageText = {
    createMeetingButton: "צור פגישה חדשה",
    joinMeetingButton: "הצטרף לפגישה",
    submitMeetingIDButton: "המשך",
    cancelMeetingIDButton: "ביטול",
    meetingIDInputPlaceholder: "הכנס מזהה פגישה",
};

export const MeetingExitText = {
    popUpTitles: {
        wrongID: "מזהה פגישה שגוי",
        error: "אירעה שגיאה",
        left: "יצאת מהפגישה",
        kicked: "הועפת מהפגישה",
        ended: "הפגישה נגמרה",
        banned: "ממש לא רוצים אותך פה",
        default: "יצאת מהפגישה",
    },
    popUpSubtitle: "מיד תועבר לדף הבית",
    popUpButton: "חזרה לדף הבית",
};

export const HostOptions = {
    kick: "הסר מהפגישה",
}