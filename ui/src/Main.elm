module Main exposing (Config, ConfigFile, Model, Msg(..), configDecoder, configEncoder, configFileEncoder, emptyConfg, expectRetrieve, expectStatus, init, main, parseDictEntry, subscriptions, update, view)

import Bootstrap.Button as Button
import Bootstrap.CDN as CDN
import Bootstrap.Card as Card
import Bootstrap.Card.Block as Block
import Bootstrap.Form as Form
import Bootstrap.Form.Input as Input
import Bootstrap.Grid as Grid
import Bootstrap.Grid.Col as Col
import Bootstrap.Utilities.Spacing as Spacing
import Browser
import Browser.Navigation as Nav
import Debug exposing (log)
import Dict
import Helpers exposing (faIcons, fontAwesome, parse)
import Html exposing (..)
import Html.Attributes exposing (..)
import Html.Events exposing (onClick)
import Http
import Json.Decode as D
import Json.Encode as E
import Url
import Url.Builder as B


main : Program () Model Msg
main =
    Browser.application
        { init = init
        , update = update
        , view = view
        , subscriptions = subscriptions
        , onUrlChange = UrlChanged
        , onUrlRequest = LinkClicked
        }


type alias ConfigFile =
    { fileName : String
    , directory : String
    }


type alias Config =
    { httpPort : String
    , thingID : String
    , thingKey : String
    , ctrlChan : String
    , dataChan : String
    , logLevel : String
    , edgexURL : String
    , mqttURL : String
    , encrypt : String
    }


emptyConfg : Config
emptyConfg =
    Config "" "" "" "" "" "" "" "" "no"


type alias Model =
    { key : Nav.Key
    , url : String
    , config : Config
    , encrypt : Bool
    , gotConfig : Bool
    , response : String
    }


type Msg
    = LinkClicked Browser.UrlRequest
    | UrlChanged Url.Url
    | PostConfig
    | GetConfig
    | GotConfig (Result Http.Error (Dict.Dict String String))
    | PostedConfig (Result Http.Error String)
    | SubmitPort String
    | SubmitThingID String
    | SubmitThingKey String
    | SubmitCtrlChan String
    | SubmitDataChan String
    | SubmitLogLevel String
    | SubmitEdgexURL String
    | SubmitMqttURL String
    | Encrypt


init : () -> Url.Url -> Nav.Key -> ( Model, Cmd Msg )
init _ url key =
    ( Model
        key
        (parse url)
        emptyConfg
        False
        False
        ""
    , Cmd.none
    )


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    let
        config =
            model.config
    in
    case msg of
        LinkClicked urlRequest ->
            case urlRequest of
                Browser.Internal url ->
                    ( model, Nav.pushUrl model.key (Url.toString url) )

                Browser.External href ->
                    ( model, Cmd.none )

        UrlChanged url ->
            ( { model | url = parse url }
            , Cmd.none
            )

        GetConfig ->
            ( model
            , Http.get
                { url =
                    B.crossOrigin
                        "http://localhost/agent/"
                        [ "config" ]
                        [ B.string "filename" "config", B.string "dir" "." ]
                , expect = expectRetrieve GotConfig configDecoder
                }
            )

        GotConfig result ->
            case result of
                Ok cfgDict ->
                    let
                        cfg =
                            Config
                                (parseDictEntry "http_port" cfgDict)
                                (parseDictEntry "mainflux_id" cfgDict)
                                (parseDictEntry "mainflux_key" cfgDict)
                                (parseDictEntry "control_channel" cfgDict)
                                (parseDictEntry "data_channel" cfgDict)
                                (parseDictEntry "log_level" cfgDict)
                                (parseDictEntry "edgex_url" cfgDict)
                                (parseDictEntry "mqtt_url" cfgDict)
                                (parseDictEntry "encrypt" cfgDict)
                    in
                    ( { model | config = cfg, gotConfig = True }, Cmd.none )

                Err error ->
                    ( model, Cmd.none )

        PostConfig ->
            ( model
            , Http.post
                { url = "http://localhost/agent/config"
                , body = Http.jsonBody <| configEncoder model.config
                , expect = expectStatus PostedConfig
                }
            )

        PostedConfig result ->
            case result of
                Ok statusCode ->
                    ( { model | response = statusCode }, Cmd.none )

                Err error ->
                    ( { model | response = Helpers.handleError error }, Cmd.none )

        SubmitPort httpPort ->
            ( { model | config = { config | httpPort = httpPort }, gotConfig = True }, Cmd.none )

        SubmitThingID thingID ->
            ( { model | config = { config | thingID = thingID }, gotConfig = True }, Cmd.none )

        SubmitThingKey thingKey ->
            ( { model | config = { config | thingKey = thingKey }, gotConfig = True }, Cmd.none )

        SubmitCtrlChan ctrlChan ->
            ( { model | config = { config | ctrlChan = ctrlChan }, gotConfig = True }, Cmd.none )

        SubmitDataChan dataChan ->
            ( { model | config = { config | dataChan = dataChan }, gotConfig = True }, Cmd.none )

        SubmitLogLevel logLevel ->
            ( { model | config = { config | logLevel = logLevel }, gotConfig = True }, Cmd.none )

        SubmitEdgexURL edgexURL ->
            ( { model | config = { config | edgexURL = edgexURL }, gotConfig = True }, Cmd.none )

        SubmitMqttURL mqttURL ->
            ( { model | config = { config | mqttURL = mqttURL }, gotConfig = True }, Cmd.none )

        Encrypt ->
            let
                encrypt =
                    not model.encrypt

                answer =
                    if encrypt then
                        "yes"

                    else
                        "no"
            in
            ( { model | encrypt = encrypt, config = { config | encrypt = answer } }, Cmd.none )


parseDictEntry : String -> Dict.Dict String String -> String
parseDictEntry key dict =
    case Dict.get key dict of
        Just value ->
            value

        Nothing ->
            ""


subscriptions : Model -> Sub Msg
subscriptions model =
    Sub.none


view : Model -> Browser.Document Msg
view model =
    { title = "Mainflux"
    , body =
        [ Grid.container []
            [ CDN.stylesheet -- creates an inline style node with the Bootstrap CSS
            , fontAwesome
            , Grid.row []
                [ Grid.col []
                    [ Card.config []
                        |> Card.headerH3 []
                            [ Grid.row []
                                [ Grid.col []
                                    [ div [ class "table_header" ]
                                        [ i [ style "margin-right" "15px", class faIcons.settings ] []
                                        , text "Config"
                                        ]
                                    ]
                                , Grid.col [ Col.attrs [ align "left" ] ]
                                    [ Form.group []
                                        [ Button.button [ Button.secondary, Button.attrs [ Spacing.ml1 ], Button.onClick GetConfig ] [ text "Get" ]
                                        ]
                                    ]
                                , Grid.col [ Col.attrs [ align "right" ] ]
                                    [ Form.group []
                                        [ Button.button [ Button.secondary, Button.attrs [ Spacing.ml1 ], Button.onClick PostConfig ] [ text "Post" ]
                                        ]
                                    ]
                                ]
                            ]
                        |> Card.block []
                            [ Block.custom
                                (Form.form []
                                    [ formGroup model.gotConfig model.config.httpPort "httpPort" "Server port" "HTTP server port" SubmitPort
                                    , formGroup model.gotConfig model.config.thingID "thingID" "Thing ID" "Mainflux thing ID" SubmitThingID
                                    , formGroup model.gotConfig model.config.thingKey "thingKey" "Thing Key" "Mainflux thing key" SubmitThingKey
                                    , formGroup model.gotConfig model.config.ctrlChan "ctrlChan" "Control channel" "Mainflux channel for sending exec and admin messages." SubmitCtrlChan
                                    , formGroup model.gotConfig model.config.dataChan "dataChan" "Data channel" "Mainflux channel for sending data messages." SubmitDataChan
                                    , formGroup model.gotConfig model.config.logLevel "logLevel" "Log level" "Mainflux logger verbosity level." SubmitLogLevel
                                    , formGroup model.gotConfig model.config.edgexURL "edgexURL" "Edgex URL" "Edgex server URL address." SubmitEdgexURL
                                    , formGroup model.gotConfig model.config.mqttURL "mqttURL" "MQTT URL" "MQTT service URL address." SubmitMqttURL
                                    , Form.group []
                                        [ input [ type_ "checkbox", onClick Encrypt ] []
                                        , Form.label [ for "encrypt" ] [ text "Encrypt" ]
                                        , Form.help [] [ text "Encrypt configuration." ]
                                        ]
                                    ]
                                )
                            ]
                        |> Card.view
                    ]
                ]
            ]
        ]
    }


formGroup : Bool -> String -> String -> String -> String -> (String -> Msg) -> Html Msg
formGroup gotConfig val id name desc msg =
    Form.group []
        [ Form.label [ for id ] [ text name ]
        , if gotConfig then
            Input.text [ Input.id id, Input.onInput msg, Input.value val ]

          else
            Input.text [ Input.id id, Input.onInput msg ]
        , Form.help [] [ text desc ]
        ]


configFileEncoder : ConfigFile -> E.Value
configFileEncoder configFile =
    E.object
        [ ( "file_name", E.string configFile.fileName )
        , ( "directory", E.string configFile.directory )
        ]


configEncoder : Config -> E.Value
configEncoder config =
    E.object
        [ ( "http_port", E.string config.httpPort )
        , ( "mainflux_id", E.string config.thingID )
        , ( "mainflux_key", E.string config.thingKey )
        , ( "control_channel", E.string config.ctrlChan )
        , ( "data_channel", E.string config.dataChan )
        , ( "log_level", E.string config.logLevel )
        , ( "edgex_url", E.string config.edgexURL )
        , ( "mqtt_url", E.string config.mqttURL )
        , ( "encrypt", E.string config.encrypt )
        ]


configDecoder : D.Decoder (Dict.Dict String String)
configDecoder =
    D.dict D.string



-- EXPECT


expectStatus : (Result Http.Error String -> msg) -> Http.Expect msg
expectStatus toMsg =
    Http.expectStringResponse toMsg <|
        \resp ->
            case resp of
                Http.BadUrl_ u ->
                    Err (Http.BadUrl u)

                Http.Timeout_ ->
                    Err Http.Timeout

                Http.NetworkError_ ->
                    Err Http.NetworkError

                Http.BadStatus_ metadata body ->
                    Err (Http.BadStatus metadata.statusCode)

                Http.GoodStatus_ metadata _ ->
                    Ok (String.fromInt metadata.statusCode)


expectRetrieve : (Result Http.Error (Dict.Dict String String) -> Msg) -> D.Decoder (Dict.Dict String String) -> Http.Expect Msg
expectRetrieve toMsg decoder =
    Http.expectStringResponse toMsg <|
        \resp ->
            case resp of
                Http.BadUrl_ u ->
                    Err (Http.BadUrl u)

                Http.Timeout_ ->
                    Err Http.Timeout

                Http.NetworkError_ ->
                    Err Http.NetworkError

                Http.BadStatus_ metadata body ->
                    Err (Http.BadStatus metadata.statusCode)

                Http.GoodStatus_ metadata body ->
                    case D.decodeString decoder body of
                        Ok value ->
                            Ok value

                        Err err ->
                            Err (Http.BadBody (D.errorToString err))
