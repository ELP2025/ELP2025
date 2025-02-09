module Test exposing (..)

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
            let
                _ = Debug.log "Message reçu : UpdateInput" newText
            in
            { model | userInput = newText }

        Draw ->
            let
                _ = Debug.log "Message reçu : Draw" model.userInput
                newParsed = divise model.userInput
                _ = Debug.log "Résultat de divise" newParsed
                actions_list = parseCommands newParsed
                _ = Debug.log "Commandes parsées" actions_list
                new_draw_list = executeCommands actions_list model
                _ = Debug.log "Dessins générés" new_draw_list
            in
            { model
                | parsedInput = model.userInput
                , angle = 0
                , position = { x = 350, y = 350 }
                , userInput = ""
                , printext = newParsed
                , drawing = new_draw_list
            }


-- PARSER
divise : String -> List String
divise input =
    let
        _ = Debug.log "Entrée utilisateur avant traitement" input
        tokens =
            input
                |> String.replace "[" " [ "
                |> String.replace "]" " ] "
                |> String.replace "," " , "
                |> String.split " "
                |> List.filter (\s -> s /= "")
        _ = Debug.log "Tokens générés" tokens
    in
    tokens


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
    let
        _ = Debug.log "Tokens reçus pour parsing" informations
    in
    case informations of
        [] ->
            let
                _ = Debug.log "Parsing terminé" []
            in
            []

        "Forward" :: n :: rest ->
            let
                _ = Debug.log "Parsing FORWARD" n
            in
            case String.toFloat n of
                Just distance ->
                    let
                        _ = Debug.log "Distance parsée" distance
                    in
                    Forward distance :: parseCommands rest
                Nothing ->
                    let
                        _ = Debug.log "Erreur de parsing FORWARD" n
                    in
                    parseCommands rest

        "Left" :: n :: rest ->
            let
                _ = Debug.log "Parsing LEFT" n
            in
            case String.toFloat n of
                Just angle ->
                    let
                        _ = Debug.log "Angle parsé" angle
                    in
                    Left angle :: parseCommands rest
                Nothing ->
                    let
                        _ = Debug.log "Erreur de parsing LEFT" n
                    in
                    parseCommands rest

        "Right" :: n :: rest ->
            let
                _ = Debug.log "Parsing RIGHT" n
            in
            case String.toFloat n of
                Just angle ->
                    let
                        _ = Debug.log "Angle parsé" angle
                    in
                    Right angle :: parseCommands rest
                Nothing ->
                    let
                        _ = Debug.log "Erreur de parsing RIGHT" n
                    in
                    parseCommands rest

        "Repeat" :: n :: "[" :: rest ->
            let
                _ = Debug.log "Parsing REPEAT" n
            in
            case String.toInt n of
                Just times ->
                    let
                        (repeatCommands, remainingTokens) = extractRepeatBlock rest [] 1
                        _ = Debug.log "Commandes répétées" repeatCommands
                        _ = Debug.log "Tokens restants" remainingTokens
                    in
                    Repeat times (parseCommands repeatCommands) :: parseCommands remainingTokens
                Nothing ->
                    let
                        _ = Debug.log "Erreur de parsing REPEAT" n
                    in
                    parseCommands rest

        _ :: rest ->
            let
                _ = Debug.log "Token non reconnu" informations
            in
            parseCommands rest


extractRepeatBlock : List String -> List String -> Int -> (List String, List String)
extractRepeatBlock informations collected bracketCount =
    let
        _ = Debug.log "Extraction du bloc REPEAT" (informations, collected, bracketCount)
    in
    case informations of
        [] ->
            let
                _ = Debug.log "Fin de l'extraction (liste vide)" collected
            in
            (collected, [])

        "[" :: rest ->
            let
                _ = Debug.log "Crochet ouvert rencontré, incrémentation de bracketCount" (bracketCount + 1)
            in
            extractRepeatBlock rest (collected ++ ["["]) (bracketCount + 1)

        "]" :: rest ->
            if bracketCount == 1 then
                let
                    _ = Debug.log "Crochet fermant correspondant trouvé, fin de l'extraction" collected
                in
                (collected, rest)
            else
                let
                    _ = Debug.log "Crochet fermant imbriqué, décrémentation de bracketCount" (bracketCount - 1)
                in
                extractRepeatBlock rest (collected ++ ["]"]) (bracketCount - 1)

        partial_information :: rest ->
            let
                _ = Debug.log "Token ajouté au bloc REPEAT" partial_information
            in
            extractRepeatBlock rest (collected ++ [partial_information]) bracketCount


executeCommands : List Command -> Model -> List (Svg Msg)
executeCommands commands model =
    let
        _ = Debug.log "Commandes à exécuter" commands
    in
    List.foldl executeCommand (model, []) commands |> Tuple.second


executeCommand : Command -> (Model, List (Svg Msg)) -> (Model, List (Svg Msg))
executeCommand command (model, drawings) =
    case command of
        Forward distance ->
            let
                newPosition = forward model.position model.angle distance
                newLine = lineSvg model.position newPosition
                _ = Debug.log "Nouvelle position après FORWARD" newPosition
            in
            ( { model | position = newPosition }, newLine :: drawings )

        Left angleChange ->
            let
                newAngle = left model.angle angleChange
                _ = Debug.log "Nouvel angle après LEFT" newAngle
            in
            ( { model | angle = newAngle }, drawings )

        Right angleChange ->
            let
                newAngle = right model.angle angleChange
                _ = Debug.log "Nouvel angle après RIGHT" newAngle
            in
            ( { model | angle = newAngle }, drawings )

        Repeat n subcommands ->
            let
                _ = Debug.log "Exécution de REPEAT" (n, subcommands)
                (updatedModel, newDrawings) =
                    List.foldl executeCommand (model, []) (List.concat (List.repeat n subcommands))
            in
            ( updatedModel, newDrawings ++ drawings )

lineSvg : Point -> Point -> Svg Msg
lineSvg p1 p2 =
    let
        _ = Debug.log "Ligne SVG générée" (p1, p2)
    in
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