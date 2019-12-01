module Helpers exposing (faIcons, fontAwesome, handleError, parse)

import Html exposing (..)
import Html.Attributes exposing (..)
import Http
import Url
import Url.Parser as UrlParser exposing ((</>))



-- URL


parse : Url.Url -> String
parse url =
    UrlParser.parse
        (UrlParser.map Tuple.pair (UrlParser.string </> UrlParser.fragment identity))
        url
        |> (\route ->
                case route of
                    Just r ->
                        Tuple.first r

                    Nothing ->
                        ""
           )



-- HTTP


handleError : Http.Error -> String
handleError error =
    case error of
        Http.BadUrl url ->
            "Bad URL: " ++ url

        Http.Timeout ->
            "Timeout"

        Http.NetworkError ->
            "Network error"

        Http.BadStatus code ->
            "Bad status: " ++ String.fromInt code

        Http.BadBody err ->
            "Invalid response body"



-- FONT AWESOME


fontAwesome : Html msg
fontAwesome =
    node "link"
        [ rel "stylesheet"
        , href "https://use.fontawesome.com/releases/v5.7.2/css/all.css"
        , attribute "integrity" "sha384-fnmOCqbTlWIlj8LyTjo7mOUStjsKC4pOpQbqyi7RrhN7udi9RwhKkMHpvLbHG9Sr"
        , attribute "crossorigin" "anonymous"
        ]
        []


faIcons =
    { provision = "fa fa-plus"
    , edit = "fa fa-pen"
    , remove = "fa fa-trash-alt"
    , dashboard = "fas fa-chart-bar"
    , things = "fas fa-sitemap"
    , channels = "fas fa-broadcast-tower"
    , connection = "fas fa-plug"
    , messages = "far fa-paper-plane"
    , version = "fa fa-code-branch"
    , websocket = "fas fa-arrows-alt-v"
    , connections = "fas fa-arrows-alt-h"
    , send = "fas fa-arrow-up"
    , receive = "fas fa-arrow-down"
    , bootstrap = "far fa-hdd"
    , settings = "fas fa-cog"
    }
