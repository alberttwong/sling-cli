import 'package:flutter/material.dart';

ClipRRect makeBlueButton(BuildContext context, String text,
    {double fontSize = 20, Function()? onPressed}) {
  var button = ClipRRect(
    borderRadius: BorderRadius.circular(4),
    child: Stack(
      children: <Widget>[
        Positioned.fill(
          child: Container(
            decoration: const BoxDecoration(
              gradient: LinearGradient(
                colors: <Color>[
                  Color(0xFF0D47A1),
                  Color(0xFF1976D2),
                  Color(0xFF42A5F5),
                ],
              ),
            ),
          ),
        ),
        TextButton(
          style: TextButton.styleFrom(
            padding: const EdgeInsets.all(16.0),
            primary: Colors.white,
            textStyle: TextStyle(fontSize: fontSize),
          ),
          onPressed: onPressed,
          child: Text(text),
        ),
      ],
    ),
  );
  return button;
}
