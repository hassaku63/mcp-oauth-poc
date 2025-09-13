package main

import (
    "encoding/json"
    "flag"
    "fmt"
    "log"
    "net/url"
    "os"
    "strings"

    "mcp-oauth-poc/internal/client/authurl"
    "mcp-oauth-poc/internal/client/callback"
    "mcp-oauth-poc/internal/client/pkce"
    "mcp-oauth-poc/internal/client/store"
    "mcp-oauth-poc/internal/client/token"
)

func main() {
    log.SetFlags(0)
    if len(os.Args) < 2 {
        usage()
        os.Exit(2)
    }
    switch os.Args[1] {
    case "new":
        cmdNew(os.Args[2:])
    case "url":
        cmdURL(os.Args[2:])
    case "complete":
        cmdComplete(os.Args[2:])
    default:
        usage()
        os.Exit(2)
    }
}

func usage() {
    fmt.Println("auth-cli commands: new | url | complete")
}

func cmdNew(args []string) {
    fs := flag.NewFlagSet("new", flag.ExitOnError)
    redirectURI := fs.String("redirect-uri", "", "redirect URI (e.g. http://127.0.0.1:53219/callback)")
    resource := fs.String("resource", "", "resource identifier (audience)")
    scope := fs.String("scope", "mcp.read mcp.write", "space-delimited scopes")
    clientID := fs.String("client-id", "mcp-cli-12345", "OAuth client_id")
    _ = fs.Parse(args)

    if *redirectURI == "" {
        log.Fatal("--redirect-uri is required")
    }
    if _, err := url.Parse(*redirectURI); err != nil {
        log.Fatalf("invalid redirect-uri: %v", err)
    }
    sess, err := pkce.NewSession(*redirectURI, *resource, *scope, *clientID)
    if err != nil {
        log.Fatal(err)
    }
    if err := store.Save(sess); err != nil {
        log.Fatal(err)
    }
    enc := json.NewEncoder(os.Stdout)
    enc.SetIndent("", "  ")
    _ = enc.Encode(sess)
}

func cmdURL(args []string) {
    fs := flag.NewFlagSet("url", flag.ExitOnError)
    sessionID := fs.String("session-id", "", "session id from 'new'")
    authorizeEndpoint := fs.String("authorize-endpoint", "", "AS authorize endpoint URL")
    openBrowser := fs.Bool("open", false, "open default browser")
    _ = fs.Parse(args)
    if *sessionID == "" || *authorizeEndpoint == "" {
        log.Fatal("--session-id and --authorize-endpoint are required")
    }
    sess, err := store.Load(*sessionID)
    if err != nil {
        log.Fatal(err)
    }
    u, err := authurl.Build(*authorizeEndpoint, sess)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(u)
    if *openBrowser {
        // best-effort: print instruction instead of trying to launch cross-platform
        fmt.Fprintln(os.Stderr, "Open the above URL in your browser.")
    }
}

func cmdComplete(args []string) {
    fs := flag.NewFlagSet("complete", flag.ExitOnError)
    sessionID := fs.String("session-id", "", "session id from 'new'")
    tokenEndpoint := fs.String("token-endpoint", "", "AS token endpoint URL")
    clientID := fs.String("client-id", "mcp-cli-12345", "OAuth client_id")
    callbackURL := fs.String("callback-url", "", "full callback URL pasted from browser")
    _ = fs.Parse(args)
    if *sessionID == "" || *tokenEndpoint == "" || *callbackURL == "" {
        log.Fatal("--session-id, --token-endpoint and --callback-url are required")
    }
    sess, err := store.Load(*sessionID)
    if err != nil {
        log.Fatal(err)
    }
    code, state, err := callback.Parse(*callbackURL)
    if err != nil {
        log.Fatal(err)
    }
    if subtleNEQ(state, sess.State) {
        log.Fatal("state mismatch")
    }
    resp, err := token.Exchange(*tokenEndpoint, token.ExchangeRequest{
        Code:         code,
        RedirectURI:  sess.RedirectURI,
        ClientID:     *clientID,
        CodeVerifier: sess.CodeVerifier,
        Resource:     sess.Resource,
    })
    if err != nil {
        log.Fatal(err)
    }
    enc := json.NewEncoder(os.Stdout)
    enc.SetIndent("", "  ")
    _ = enc.Encode(resp)
    _ = store.Delete(*sessionID)
}

func subtleNEQ(a, b string) bool {
    if len(a) != len(b) {
        return true
    }
    var v byte
    for i := 0; i < len(a); i++ {
        v |= a[i] ^ b[i]
    }
    return v != 0
}

