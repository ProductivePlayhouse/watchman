import React from 'react';
import { StyleSheetManager } from 'styled-components';

const StylesProvider = ({ children }) =>
{
    const nonce = window.cspNonce === "__CSP_NONCE__" ? null : window.cspNonce;

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