module Main exposing (..)

import Browser
import Svg exposing (..)
import Svg.Attributes exposing (..)
import Html exposing (Html, button, div, input, text)
import Html.Events exposing (onClick, onInput)
import Html.Attributes as HtmlA exposing (placeholder, value)
import Json.Decode exposing (int)


-- MAIN
main : Program () (Model Msg) Msg
main =
    Browser.sandbox { init = initialModel, update = update, view = view }


-- MODEL
type alias Model msg =
    { userInput : String
    , parsedInput : String
    , printext : List String
    , drawing : List (Svg msg)
    , angle : Float
    , position : Point
    }


initialModel : Model Msg
initialModel =
    { userInput = ""
    , parsedInput = ""
    , printext = []
    , drawing = []
    , angle = 0
    , position = { x = 350, y = 350 } -- Position initiale de la tortue

    }


-- MESSAGES
type Msg
    = UpdateInput String
    | Draw


-- COMMANDES
type Command
    = Forward Float
    | Left Float
    | Right Float
    | Repeat Int (List Command)


type alias Point = { x : Float, y : Float}


-- UPDATE
update : Msg -> Model Msg -> Model Msg
update msg model =
    case msg of
        UpdateInput newText ->
            { model | userInput = newText }

        Draw ->
            let
                newParsed = divise model.userInput
            in
            { model
                | parsedInput = model.userInput
                , userInput = ""
                , printext = newParsed
                , drawing = []
                , angle = 0
            }


-- PARSER
divise : String -> List String
divise input =
    input
        |> String.replace "[" " [ "
        |> String.replace "]" " ] "
        |> String.split " "
        |> List.filter (\s -> s /= "")

-- Fonction RIGHT
right : Float -> Float -> Float
right angle input_angle = 
    angle + input_angle

-- Fonction LEFT
left : Float -> Float -> Float
left angle input_angle = 
    angle - input_angle

-- Fonction FORWARD
forward : Point -> Float -> Float -> Point
forward {x, y} angle input_avance =
    { x = x + input_avance * cos (degrees angle)
    , y = y + input_avance * sin (degrees angle)}

-- Fonction REPEAT
repeat : Int -> List (Point -> Point) -> Point -> Point
repeat n actions model =
    List.foldl (\_ m -> applyActions actions m) model (List.repeat n ())

-- Appliquer une liste d'actions à un modèle
applyActions : List (Point -> Point) -> Point -> Point
applyActions actions model =
    List.foldl (\action m -> action m) model actions
 


-- VIEW
view : Model Msg -> Html Msg
view model =
    div
        [ HtmlA.style "display" "flex"
        , HtmlA.style "flex-direction" "column"
        , HtmlA.style "align-items" "center"
        , HtmlA.style "margin-top" "20px"
        ]
        [ div
            [ HtmlA.style "padding" "10px"
            , HtmlA.style "width" "100%"
            , HtmlA.style "text-align" "center"
            ]
            [ Html.text "Type in your code below" ]

        , input
            [ placeholder "Écrivez ici..."
            , value model.userInput
            , onInput UpdateInput
            , HtmlA.style "margin" "10px"
            , HtmlA.style "padding" "5px"
            , HtmlA.style "width" "45%"
            ]
            []

        , button
            [ onClick Draw
            , HtmlA.style "padding" "10px"
            , HtmlA.style "margin-top" "10px"
            , HtmlA.style "width" "150px"
            ]
            [ Html.text "Draw" ]

        , svg
            [ width "700"
            , height "700"
            , viewBox "0 0 120 120"
            , HtmlA.style "margin-top" "10px"
            , HtmlA.style "border" "2px solid black"
            , HtmlA.style "background-color" "white"
            ]
            model.drawing  -- Correction de l'affichage du dessin

        , div
            [ HtmlA.style "color" "red"
            , HtmlA.style "margin-top" "20px"
            ]
            [ Html.text (String.join " " model.printext) ]
        ]
