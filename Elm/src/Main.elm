module Main exposing (..)

import Browser
import Svg exposing (..)
import Svg.Attributes exposing (..)
import Html exposing (Html, button, div, input)
import Html.Events exposing (onClick, onInput)
import Html.Attributes as HtmlA exposing (placeholder, value)

-- MAIN
main = Browser.sandbox { init = initialModel, update = update, view = view }

-- MODEL
type alias Model = { userInput : String
                    , parsedInput : String
                    , printext : List String}

initialModel : Model
initialModel = { userInput = "" 
                , parsedInput = ""
                , printext = []}

-- MESSAGE (Msg)
type Msg = UpdateInput String | Draw

--Commande Tcturtle
type Command
    = Forward Int
    | Left Int
    | Right Int
    | Repeat Int (List Command)


-- UPDATE
update : Msg -> Model -> Model
update msg model =
    case msg of
        UpdateInput newText -> { model | userInput = newText }
        Draw -> { model | parsedInput = model.userInput, userInput = "", printext = divise(model.parsedInput)}

-- FONCTION DE DÉCOUPAGE DES INSTRUCTIONS
divise : String -> List String
divise input =
    input
        |> String.replace "[" " [ "
        |> String.replace "]" " ] "
        |> String.split " "
        |> List.filter (\s -> s /= "")


-- VIEW
view : Model -> Html Msg
view model =
    div [ HtmlA.style "display" "flex", HtmlA.style "flex-direction" "column", HtmlA.style "align-items" "center", HtmlA.style "margin-top" "20px" ]
        [ div [ HtmlA.style "padding" "10px", HtmlA.style "width" "100%", HtmlA.style "text-align" "center" ] [ Html.text "Type in your code below" ]
        , input [ placeholder "Écrivez ici...", value model.userInput, onInput UpdateInput, HtmlA.style "margin" "10px", HtmlA.style "padding" "5px", HtmlA.style "width" "45%" ] []
        , button [ onClick Draw, HtmlA.style "padding" "10px", HtmlA.style "margin-top" "10px", HtmlA.style "width" "150px" ] [ Html.text "Draw" ]
        , svg
            [ width "700", height "700", viewBox "0 0 120 120"
            , HtmlA.style "margin-top" "10px", HtmlA.style "border" "2px solid black", HtmlA.style "background-color" "white"
            ]
            [ line [ x1 "10", y1 "10", x2 "100", y2 "10", stroke "black", strokeWidth "2" ] [] 
            , line [ x1 "30", y1 "20", x2 "200", y2 "20", stroke "black", strokeWidth "2" ] [] ]
        , div [ HtmlA.style "color" "red", HtmlA.style "margin-top" "20px" ] [ Html.text (String.join " " model.printext)]
        ]
