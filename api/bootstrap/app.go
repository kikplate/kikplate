package bootstrap

import (
	"github.com/kickplate/api/command"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "clean-gin",
	Short: "Clean architecture using gin framework",
	Long: `

                                                                                                                                      
                                                                                                                                      
kkkkkkkk             iiii  kkkkkkkk                               lllllll                           tttt                              
k::::::k            i::::i k::::::k                               l:::::l                        ttt:::t                              
k::::::k             iiii  k::::::k                               l:::::l                        t:::::t                              
k::::::k                   k::::::k                               l:::::l                        t:::::t                              
 k:::::k    kkkkkkkiiiiiii  k:::::k    kkkkkkkppppp   ppppppppp    l::::l   aaaaaaaaaaaaa  ttttttt:::::ttttttt        eeeeeeeeeeee    
 k:::::k   k:::::k i:::::i  k:::::k   k:::::k p::::ppp:::::::::p   l::::l   a::::::::::::a t:::::::::::::::::t      ee::::::::::::ee  
 k:::::k  k:::::k   i::::i  k:::::k  k:::::k  p:::::::::::::::::p  l::::l   aaaaaaaaa:::::at:::::::::::::::::t     e::::::eeeee:::::ee
 k:::::k k:::::k    i::::i  k:::::k k:::::k   pp::::::ppppp::::::p l::::l            a::::atttttt:::::::tttttt    e::::::e     e:::::e
 k::::::k:::::k     i::::i  k::::::k:::::k     p:::::p     p:::::p l::::l     aaaaaaa:::::a      t:::::t          e:::::::eeeee::::::e
 k:::::::::::k      i::::i  k:::::::::::k      p:::::p     p:::::p l::::l   aa::::::::::::a      t:::::t          e:::::::::::::::::e 
 k:::::::::::k      i::::i  k:::::::::::k      p:::::p     p:::::p l::::l  a::::aaaa::::::a      t:::::t          e::::::eeeeeeeeeee  
 k::::::k:::::k     i::::i  k::::::k:::::k     p:::::p    p::::::p l::::l a::::a    a:::::a      t:::::t    tttttte:::::::e           
k::::::k k:::::k   i::::::ik::::::k k:::::k    p:::::ppppp:::::::pl::::::la::::a    a:::::a      t::::::tttt:::::te::::::::e          
k::::::k  k:::::k  i::::::ik::::::k  k:::::k   p::::::::::::::::p l::::::la:::::aaaa::::::a      tt::::::::::::::t e::::::::eeeeeeee  
k::::::k   k:::::k i::::::ik::::::k   k:::::k  p::::::::::::::pp  l::::::l a::::::::::aa:::a       tt:::::::::::tt  ee:::::::::::::e  
kkkkkkkk    kkkkkkkiiiiiiiikkkkkkkk    kkkkkkk p::::::pppppppp    llllllll  aaaaaaaaaa  aaaa         ttttttttttt      eeeeeeeeeeeeee  
                                               p:::::p                                                                                
                                               p:::::p                                                                                
                                              p:::::::p                                                                               
                                              p:::::::p                                                                               
                                              p:::::::p                                                                               
                                              ppppppppp                                                                               
                                                                                                                                      
                                         		
This is a command runner or cli for api architecture in golang. 
Using this we can use underlying dependency injection container for running scripts. 
Main advantage is that, we can use same services, repositories, infrastructure present in the application itself`,
	TraverseChildren: true,
}

type App struct {
	*cobra.Command
}

func NewApp() App {
	cmd := App{
		Command: rootCmd,
	}
	cmd.AddCommand(command.GetSubCommands(CommonModules)...)
	return cmd
}

var RootApp = NewApp()
