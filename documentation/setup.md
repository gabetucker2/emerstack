# How to Set Up Emerstack

## Step 1) Download the prerequisite software

> Download the programming language [Go](https://go.dev/doc/install) so that you can run Go code

> Download [Git](https://git-scm.com/downloads) so that you can clone Emerstack to your computer

## Step 2) Make a local copy of Emergent

*[Leabra](https://en.wikipedia.org/wiki/Leabra) is a biologically plausible neural network algorithm.  [Dr. Randy O'Reilly](https://en.wikipedia.org/wiki/Randall_C._O%27Reilly) developed an open-source implementation of this algorithm called [Emergent](https://en.wikipedia.org/wiki/Emergent_(software)) that we will be using to implement Emerstack.*

> A) Open a [command prompt window](https://www.lifewire.com/how-to-open-command-prompt-2618089)

> B) Navigate to the directory in which you would like to store Emerstack
>> Navigating directories with a [Mac](https://techwiser.com/how-to-navigate-to-a-folder-in-terminal-mac/#:~:text=1%20Method%20I.%20This%20is%20the%20most%20usual,to%20navigate%20to%20a%20folder%20in%20the%20terminal.) PC
>
>> Navigating directories with a [Windows](https://techwiser.com/how-to-navigate-to-a-folder-in-terminal-mac/#:~:text=1%20Method%20I.%20This%20is%20the%20most%20usual,to%20navigate%20to%20a%20folder%20in%20the%20terminal.) PC

> C) Clone Emergent to your computer by running the following command:
>> ```
>> git clone https://github.com/emer/leabra
>> ```

## Step 3) Make a local copy of Emerstack

> A) After leabra finishes cloning to your designated directory, navigate to ```leabra/examples``` with the following command:
>> ```
>> cd leabra/examples
>> ```

> B) Clone Emerstack to your computer by running the following command:
>> ```
>> git clone https://github.com/gabetucker2/emerstack
>> ```

## Step 4) Ensure go.mod is up-to-date

*go.mod is a file that keeps track of import versions so that you can manage your versions all in one place, rather than having to type v1.2.3 in each individual script's import.*

*Since we imported Emerstack inside of our Leabra directory, Emerstack is dependent on our Leabra directory's go.mod file.  Therefore, we must run ```go mod tidy``` in our Leabra directory to import the appropriate Emerstack dependencies to its parent directory, Leabra.*

> A) Return to the ```leabra``` directory in your terminal by running the following command:
>> ```
>> cd ../
>> ```

> B) Import Emerstack dependencies by running the following command:
>> ```
>> go mod tidy
>> ```

## Step 5) Configure and run Emerstack

*Check out our other tutorials to learn how to configure Emerstack.  Once you have configured it to your liking, do the following:*

> A) Return to the ```leabra/examples/emerstack``` in your terminal by running the following command:
>> ```
>> cd examples/leabra
>> ```

> B) Build Emerstack by running the following command:
>> ```
>> go run .
>> ```

> C) After Emergent's user interface launches, you can play with the model to your liking!
