module Test exposing (..)

import Browser
import Svg exposing (..)
import Svg.Attributes exposing (..)
import Html exposing (Html, button, div, input, text)
import Html.Events exposing (onClick, onInput)
import Html.Attributes as HtmlA exposing (placeholder, value)
import Json.Decode exposing (int)


-- MAIN
main : Program () Model Msg
main =
    Browser.sandbox { init = initialModel, update = update, view = view }


-- MODEL
type alias Model =
    { userInput : String
    , parsedInput : String
    , printext : List String
    , drawing : List (Svg Msg)
    , angle : Float
    , position : Point
    }


initialModel : Model
initialModel =
    { userInput = ""
    , parsedInput = ""
    , printext = []
    , drawing = []
    , angle = 0
    , position = { x = 60, y = 60 } -- Position initiale de la tortue
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
update : Msg -> Model -> Model
update msg model =
    case msg of
        UpdateInput newText ->
            { model | userInput = newText }

        Draw ->
            let
                newParsed = divise model.userInput
                commands = parseCommands newParsed
                newDrawing = executeCommands commands model.position model.angle
            in
            { model
                | parsedInput = model.userInput
                , userInput = ""
                , printext = newParsed
                , drawing = newDrawing
                , angle = 0
                , position = { x = 350, y = 350 }
            }


-- PARSER
divise : String -> List String
divise input =
    input
        |> String.replace "[" " [ "
        |> String.replace "]" " ] "
        |> String.split " "
        |> List.filter (\s -> s /= "")


-- Parseur des commandes
parseCommands : List String -> List Command
parseCommands [] = []
parseCommands ("Repeat" :: nStr :: rest) =
    case String.toInt nStr of
        Just n ->
            let
                (commands, remaining) = parseRepeatCommands rest
            in
            Repeat n commands :: parseCommands remaining
        Nothing -> parseCommands rest
parseCommands ("Left" :: angleStr :: rest) =
    case String.toFloat angleStr of
        Just angle -> Left angle :: parseCommands rest
        Nothing -> parseCommands rest
parseCommands ("Right" :: angleStr :: rest) =
    case String.toFloat angleStr of
        Just angle -> Right angle :: parseCommands rest
        Nothing -> parseCommands rest
parseCommands ("Forward" :: distanceStr :: rest) =
    case String.toFloat distanceStr of
        Just distance -> Forward distance :: parseCommands rest
        Nothing -> parseCommands rest
parseCommands _ = []


-- Parse les commandes "Repeat"
parseRepeatCommands : List String -> (List Command, List String)
parseRepeatCommands [] = ([], [])
parseRepeatCommands ("]" :: rest) = ([], rest)  -- Fin de la répétition
parseRepeatCommands ("Left" :: angleStr :: rest) =
    case String.toFloat angleStr of
        Just angle ->
            let (commands, remaining) = parseRepeatCommands rest
            in (Left angle :: commands, remaining)
        Nothing -> parseRepeatCommands rest
parseRepeatCommands ("Right" :: angleStr :: rest) =
    case String.toFloat angleStr of
        Just angle ->
            let (commands, remaining) = parseRepeatCommands rest
            in (Right angle :: commands, remaining)
        Nothing -> parseRepeatCommands rest
parseRepeatCommands ("Forward" :: distanceStr :: rest) =
    case String.toFloat distanceStr of
        Just distance ->
            let (commands, remaining) = parseRepeatCommands rest
            in (Forward distance :: commands, remaining)
        Nothing -> parseRepeatCommands rest
parseRepeatCommands _ = ([], [])


-- Exécuter les commandes
executeCommands : List Command -> Point -> Float -> List (Svg Msg)
executeCommands commands position angle =
    List.foldl (\cmd acc -> executeCommand cmd position angle acc) [] commands


-- Exécuter une seule commande
executeCommand : Command -> Point -> Float -> List (Svg Msg) -> List (Svg Msg)
executeCommand cmd position angle acc =
    case cmd of
        Forward distance ->
            let
                newPosition = forward position angle distance
                newLine = line [ x1 (String.fromFloat position.x), y1 (String.fromFloat position.y), x2 (String.fromFloat newPosition.x), y2 (String.fromFloat newPosition.y), stroke "black", strokeWidth "2" ] []
            in
            newLine :: acc
        Left angleChange ->
            Left (angle + angleChange) :: acc
        Right angleChange ->
            Right (angle - angleChange) :: acc
        Repeat n subCommands ->
            let
                subDrawing = executeCommands subCommands position angle
            in
            List.repeat n subDrawing |> List.concat |> List.append acc


-- FORWARD
forward : Point -> Float -> Float -> Point
forward {x, y} angle input_avance =
    { x = x + input_avance * cos (degrees angle)
    , y = y + input_avance * sin (degrees angle)
    }


-- VIEW
view : Model -> Html Msg
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
