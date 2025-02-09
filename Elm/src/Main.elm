module Main exposing (..)

import Browser
import Svg exposing (..)
import Svg.Attributes exposing (..)
import Html exposing (Html, button, div, input)
import Html.Events exposing (onClick, onInput)
import Html.Attributes as HtmlA exposing (placeholder, value)


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


type alias Point = { x : Float, y : Float }


-- UPDATE
update : Msg -> Model -> Model
update msg model =
    case msg of
        UpdateInput newText ->
            { model | userInput = newText }

        Draw ->
            let
                newParsed = divise model.userInput
                actions_list = parseCommands newParsed
                new_draw_list = executeCommands actions_list model
            in
            { model
                | parsedInput = model.userInput
                , angle = 0
                , position = { x = 350, y = 350 } -- Position initiale de la tortue
                , userInput = ""
                , printext = newParsed
                , drawing = new_draw_list
            }


-- PARSER
divise : String -> List String
divise input =
    input
        |> String.replace "[" " [ "
        |> String.replace "]" " ] "
        |> String.replace "," " , "
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
forward { x, y } angle input_avance =
    { x = x + input_avance * cos (degrees angle)
    , y = y + input_avance * sin (degrees angle)
    }


parseCommands : List String -> List Command
parseCommands informations =
    case informations of
        [] ->
            []

        "Forward" :: n :: rest ->
            case String.toFloat n of
                Just distance -> Forward distance :: parseCommands rest
                Nothing -> parseCommands rest -- Ignore en cas d'erreur

        "Left" :: n :: rest ->
            case String.toFloat n of
                Just angle -> Left angle :: parseCommands rest
                Nothing -> parseCommands rest

        "Right" :: n :: rest ->
            case String.toFloat n of
                Just angle -> Right angle :: parseCommands rest
                Nothing -> parseCommands rest

        "Repeat" :: n :: "[" :: rest ->
            case String.toInt n of
                Just times ->
                    let
                        (repeatCommands, remainingTokens) = extractRepeatBlock rest [] 1
                    in
                    Repeat times (parseCommands repeatCommands) :: parseCommands remainingTokens

                Nothing ->
                    parseCommands rest -- Ignore en cas d'erreur

        _ :: rest ->
            parseCommands rest -- Ignore les éléments non reconnus


extractRepeatBlock : List String -> List String -> Int -> (List String, List String)
extractRepeatBlock informations collected bracketCount =
    case informations of
        [] ->
            -- Si on arrive à la fin de la liste sans fermer tous les crochets
            (collected, [])

        "[" :: rest ->
            -- On rencontre un crochet ouvert, on incrémente le compteur
            extractRepeatBlock rest (collected ++ ["["]) (bracketCount + 1)

        "]" :: rest ->
            if bracketCount == 1 then
                -- Si le compteur est à 1, on a trouvé le crochet fermant correspondant
                (collected, rest)
            else
                -- Sinon, on décrémente le compteur et on continue
                extractRepeatBlock rest (collected ++ ["]"]) (bracketCount - 1)

        partial_information :: rest ->
            -- On ajoute le token à la liste collectée
            extractRepeatBlock rest (collected ++ [partial_information]) bracketCount


executeCommands : List Command -> Model -> List (Svg Msg)
executeCommands commands model =
    List.foldl executeCommand (model, []) commands |> Tuple.second


executeCommand : Command -> (Model, List (Svg Msg)) -> (Model, List (Svg Msg))
executeCommand command (model, drawings) =
    case command of
        Forward distance ->
            let
                newPosition = forward model.position model.angle distance
                newLine = lineSvg model.position newPosition
            in
            ( { model | position = newPosition }, newLine :: drawings )

        Left angleChange ->
            ( { model | angle = left model.angle angleChange }, drawings )

        Right angleChange ->
            ( { model | angle = right model.angle angleChange }, drawings )

        Repeat n subcommands ->
            let
                (updatedModel, newDrawings) =
                    List.foldl executeCommand (model, []) (List.concat (List.repeat n subcommands))
            in
            ( updatedModel, newDrawings ++ drawings )


lineSvg : Point -> Point -> Svg Msg
lineSvg p1 p2 =
    Svg.line
        [ x1 (String.fromFloat p1.x)
        , y1 (String.fromFloat p1.y)
        , x2 (String.fromFloat p2.x)
        , y2 (String.fromFloat p2.y)
        , stroke "black"
        , strokeWidth "1"
        ]
        []


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
            , viewBox "0 0 800 800"
            , HtmlA.style "margin-top" "10px"
            , HtmlA.style "border" "2px solid black"
            , HtmlA.style "background-color" "white"
            ]
            model.drawing -- Correction de l'affichage du dessin

        , div
            [ HtmlA.style "color" "red"
            , HtmlA.style "margin-top" "20px"
            ]
            [ Html.text (String.join " " model.printext) ]
        ]