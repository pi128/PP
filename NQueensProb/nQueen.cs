using System;
using System.Collections.Generic;


namespace NQueensProb
{
    class NQueen
    {
        // Store every full 2D-board solution
        static List<int[,]> solutions = new List<int[,]>();

        public static void Main(string[] args)
        {
            Console.WriteLine("N-Queens: place N queens on an NxN board so none attack each other.");

            Console.Write("What size board anything over 12 will probably fail its recursive : ");
            int N = int.Parse(Console.ReadLine());

            int[,] board = new int[N, N]; 
            Solve(board, 0);

            Console.WriteLine($"Found {solutions.Count} solutions.\n");

            for (int i = 0; i < solutions.Count; i++)
            {
                Console.WriteLine($"Solution {i + 1}:");
                PrintBoard(solutions[i]);
                Console.WriteLine();
            }

            // If you want just one random solution instead:
            // var rng = new Random();
            // PrintBoard(solutions[rng.Next(solutions.Count)]);
        }

        // recursively solving for the boards and then also storing all solutions
        public static void Solve(int[,] board, int row)
        {
            int n = board.GetLength(0);

            // Reached a complete placement
            if (row == n)
            {
                int[,] copy = new int[n, n];
                for (int r = 0; r < n; r++)
                    for (int c = 0; c < n; c++)
                        copy[r, c] = board[r, c];

                solutions.Add(copy);
                return; // keep searching for more
            }

            for (int col = 0; col < n; col++)
            {
                if (IsSafe(board, row, col))
                {
                    board[row, col] = 1;     // place
                    Solve(board, row + 1);   // recurse
                    board[row, col] = 0;     // backtrack
                }
            }
        }

        // Safety check on the 2D board
        public static bool IsSafe(int[,] board, int row, int col)
        {
            int n = board.GetLength(0);

            // Row
            for (int c = 0; c < n; c++)
                if (board[row, c] == 1) return false;

            // Column
            for (int r = 0; r < n; r++)
                if (board[r, col] == 1) return false;

            // Top Left to Bottom Right diagonal
            for (int r = row, c = col; r >= 0 && c >= 0; r--, c--)
                if (board[r, c] == 1) return false;
            for (int r = row, c = col; r < n && c < n; r++, c++)
                if (board[r, c] == 1) return false;

            // Top Right to Bottom Left diagonal
            for (int r = row, c = col; r >= 0 && c < n; r--, c++)
                if (board[r, c] == 1) return false;
            for (int r = row, c = col; r < n && c >= 0; r++, c--)
                if (board[r, c] == 1) return false;

            return true;
        }

        // a simple 2d list printer
        public static void PrintBoard(int[,] board)
        {
            int n = board.GetLength(0);
            for (int r = 0; r < n; r++)
            {
                for (int c = 0; c < n; c++)
                    Console.Write(board[r, c] + " ");
                Console.WriteLine();
            }
        }
    }
}