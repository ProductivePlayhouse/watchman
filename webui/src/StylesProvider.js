import React from 'react';
import { StyleSheetManager } from 'styled-components';

const StylesProvider = ({ children }) =>
{
    const nonce = window.cspNonce === "__CSP_NONCE__" ? null : window.cspNonce;
    console.log("CSP Nonce: ", nonce);

    return nonce ? (
        <StyleSheetManager nonce={nonce}>
            {children}
        </StyleSheetManager>
    ) : (
        <>
            {children}
        </>
    );
};

export default StylesProvider;